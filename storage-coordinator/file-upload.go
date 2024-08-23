package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
	"github.com/lielalmog/SkyArchive/storage-coordinator/database"
	"github.com/segmentio/kafka-go"
)

func combineChunksAndUploadToStorage(fileName, destContainerName string) error {
	connectionString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return err
	}

	blobServiceClient, err := blob.NewClientFromConnectionString(connectionString, tempContainerName, fileName, nil)
	if err != nil {
		return err
	}

	blobProperties, err := blobServiceClient.GetProperties(context.Background(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Get the size of the file for 1MB in bytes
	chuckSize := 1024 * 1024
	fileSize := *blobProperties.ContentLength
	var offset int64 = 0

	reader, writer := io.Pipe()

	go func() {
		defer reader.Close()

		_, err = serviceClient.UploadStream(context.Background(), destContainerName, fileName, reader, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for offset < fileSize {
		var count int

		if offset+int64(chuckSize) > fileSize {
			count = int(fileSize - offset)
		} else {
			count = chuckSize
		}

		// This function reads from the permanent container and writes to the writer pipe
		blobDownloadResponse, err := serviceClient.DownloadStream(context.Background(), tempContainerName, fileName, &azblob.DownloadStreamOptions{
			Range: blob.HTTPRange{
				Offset: int64(offset),
				Count:  int64(count),
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		bodyStream := blobDownloadResponse.Body
		_, err = io.Copy(writer, bodyStream)
		if err != nil {
			return err
		}

		bodyStream.Close()
		offset += int64(chuckSize)
	}

	writer.Close()

	return nil
}

func fileUpload(connString string) {
	db := database.GetDB()
	validator := configs.GetValidator()

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
			wg.Done()
		}()

		go func() {
			combineChunksAndUploadToStorage(fileName, backupContainerName)
			wg.Done()
		}()

		url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", stgAccountName, permanentContainerName, fileName)
		_, err = db.Pool.Exec(context.Background(), "UPDATE files SET status = 'UPLOADED', url = $1 WHERE file_id = $2", url, *payload.FileId)

		wg.Wait()
		r.CommitMessages(context.Background(), m)

		if err != nil {
			panic(err)
		}
	}
}
