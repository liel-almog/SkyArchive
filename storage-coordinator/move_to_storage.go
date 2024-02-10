package main

import (
	"context"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
)

func combineChunksAndUploadToStorage(fileName string, destContainerName string) error {
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
