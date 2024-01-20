package main

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/lielalmog/file-uploader/storage-coordinator/configs"
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

		// bytes to int64
		id, err := strconv.ParseInt(string(m.Value), 0, 64)
		if err != nil {
			continue
		}

		wg.Add(2)
		go func() {
			combineChunksAndUploadToStorage(id, permanentContainerName)
			w.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(strconv.FormatInt(id, 10)),
				Value: []byte("permanent"),
			})
			wg.Done()
		}()

		go func() {
			combineChunksAndUploadToStorage(id, backupContainerName)
			w.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(strconv.FormatInt(id, 10)),
				Value: []byte("backup"),
			})
			wg.Done()
		}()

		// commit the message that was just read
		wg.Wait()
		r.CommitMessages(context.Background(), m)
	}
}
