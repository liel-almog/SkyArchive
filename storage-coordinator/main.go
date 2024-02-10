package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
	"github.com/segmentio/kafka-go"
)

const (
	fileUplaodFinilizationTopic  = "file-upload-finalization"
	fileUploadStatusUpdatesTopic = "file-upload-status-updates"

	permanentContainerName = "permanent-files"
	backupContainerName    = "backup-files"
	tempContainerName      = "temp-files"
)

func main() {
	configs.InitEnv()

	validator := configs.GetValidator()

	var wg sync.WaitGroup

	brokers, err := configs.GetEnv("KAFKA_BROKERS")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   fileUplaodFinilizationTopic,
		GroupID: "upload-permanent-backup",
	})

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  strings.Split(brokers, ","),
		Topic:    fileUploadStatusUpdatesTopic,
		Balancer: &kafka.LeastBytes{},
	})

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}

		payload := new(FileUploadFinalizationMessage)

		if err = json.Unmarshal(m.Value, payload); err != nil {
			panic(err)
		}

		if err = validator.Struct(payload); err != nil {
			panic(err)
		}

		fileName := strconv.FormatInt(*payload.FileId, 10)

		wg.Add(2)
		go func() {
			combineChunksAndUploadToStorage(fileName, permanentContainerName)
			w.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(fileName),
				Value: []byte("permanent"),
			})
			wg.Done()
		}()

		go func() {
			combineChunksAndUploadToStorage(fileName, backupContainerName)
			w.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(fileName),
				Value: []byte("backup"),
			})
			wg.Done()
		}()

		// commit the message that was just read
		wg.Wait()
		// r.CommitMessages(context.Background(), m)
	}
}
