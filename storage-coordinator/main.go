package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/lielalmog/file-uploader/storage-coordinator/configs"
	"github.com/segmentio/kafka-go"
)

func main() {
	configs.InitEnv()
	const topic = "file-upload-finalization"

	brokers, err := configs.GetEnv("KAFKA_BROKERS")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		GroupID: "upload-permanent-backup",
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

		go combineChunksAndUploadToPermanent(id)
	}
}
