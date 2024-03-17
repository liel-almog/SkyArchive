package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
	"github.com/segmentio/kafka-go"
)

func deleteFromStorage(connString, fileName, containerName string) error {
	serviceClient, err := azblob.NewClientFromConnectionString(connString, nil)
	if err != nil {
		return err
	}

	_, err = serviceClient.DeleteBlob(context.Background(), containerName, fileName, nil)
	if err != nil {
		return err
	}

	return nil
}

func fileDelete(connString string) {
	validator := configs.GetValidator()
	var wg sync.WaitGroup

	brokers, err := configs.GetEnv("KAFKA_BROKERS")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   fileDeleteTopic,
		GroupID: "delete-permanent-backup",
	})

	for {
		m, err := r.FetchMessage(context.Background())

		if err != nil {
			break
		}

		payload := new(FileDeleteMessage)

		if err = json.Unmarshal(m.Value, payload); err != nil {
			panic(err)
		}

		if err := validator.Struct(payload); err != nil {
			panic(err)
		}

		fileName := strconv.FormatInt(*payload.FileId, 10)

		wg.Add(2)

		go func() {
			deleteFromStorage(connString, fileName, permanentContainerName)
			wg.Done()
		}()

		go func() {
			deleteFromStorage(connString, fileName, backupContainerName)
			wg.Done()
		}()

		wg.Wait()
		err = r.CommitMessages(context.Background(), m)

		if err != nil {
			panic(err)
		}
	}
}
