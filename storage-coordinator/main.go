package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
	"github.com/lielalmog/SkyArchive/storage-coordinator/database"
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
	db := database.GetDB()

	connString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		panic(err)
	}

	var stgAccountName string

	connStringParts := strings.Split(connString, ";")
	for _, part := range connStringParts {
		if strings.Contains(part, "AccountName") {
			stgAccountName = strings.Split(part, "=")[1]
		}
	}

	if stgAccountName == "" {
		panic("Storage account name not found in connection string")
	}

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

		wg.Wait()
		r.CommitMessages(context.Background(), m)
		url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", stgAccountName, permanentContainerName, fileName)
		if err != nil {
			panic(err)
		}

		_, err = db.Pool.Exec(context.Background(), "UPDATE files SET status = 'uploaded', url = $1 WHERE file_id = $2", url, *payload.FileId)

		if err != nil {
			panic(err)
		}
	}
}
