package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/lielalmog/file-uploader/storage-coordinator/configs"
)

func combineChunksAndUploadToStorage(id int64, containerName string) error {
	connectionString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()

	blobPrefix := fmt.Sprintf("%d/", id)
	pager := serviceClient.NewListBlobsFlatPager(tempContainerName, &azblob.ListBlobsFlatOptions{
		Prefix: &blobPrefix,
	})

	// This function reads from the reader pipe and uploads the data to the permanent container as a stream
	go func() {
		defer reader.Close()

		_, err = serviceClient.UploadStream(context.Background(), containerName, fmt.Sprintf("%d", id), reader, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var chunks []string

	for pager.More() {
		// advance to the next page
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return err
		}

		for _, blob := range page.Segment.BlobItems {
			// Downloads the chunk from the temporary container and writes it to the pipe
			blobName := *blob.Name
			chunks = append(chunks, blobName)

		}
	}

	slices.Sort(chunks)

	for _, chunk := range chunks {
		blobDownloadResponse, err := serviceClient.DownloadStream(context.Background(), tempContainerName, chunk, nil)
		if err != nil {
			return err
		}

		bodyStream := blobDownloadResponse.Body
		_, err = io.Copy(writer, bodyStream)
		if err != nil {
			return err
		}

		bodyStream.Close()
	}

	writer.Close()

	return nil
}
