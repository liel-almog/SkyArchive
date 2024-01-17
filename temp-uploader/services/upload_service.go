package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/lielalmog/be-file-streaming/models"
	"github.com/lielalmog/be-file-streaming/repositories"
)

type UploadService interface {
	StartUpload(fileMetadate *models.FileMetadateDTO) (*int64, error)
	UploadChunk(fileHeader *multipart.FileHeader, id int64, chunkIndex int) error
	CombineChunksAndUploadToPermanent(id int64) error
}

type uploadServiceImpl struct{}

var (
	initUploadServiceOnce sync.Once
	uploadService         *uploadServiceImpl
)

const (
	CONTAINER_NAME           = "files"
	PERMANENT_CONTAINER_NAME = "permanent-files"
)

func (u *uploadServiceImpl) StartUpload(fileMetadate *models.FileMetadateDTO) (*int64, error) {
	return repositories.GetUploadRepository().SaveFileMetadata(fileMetadate)
}

func (u *uploadServiceImpl) UploadChunk(fileHeader *multipart.FileHeader, id int64, chunkIndex int) error {
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	blobName := fmt.Sprintf("%d/%d", id, chunkIndex)
	serviceClient.UploadStream(context.Background(), CONTAINER_NAME, blobName, file, nil)
	if err != nil {
		return err
	}

	return nil
}

func (u *uploadServiceImpl) CombineChunksAndUploadToPermanent(id int64) error {
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()

	blobPrefix := fmt.Sprintf("%d/", id)
	pager := serviceClient.NewListBlobsFlatPager(CONTAINER_NAME, &azblob.ListBlobsFlatOptions{
		Prefix: &blobPrefix,
	})

	// This function reads from the reader pipe and uploads the data to the permanent container as a stream
	go func() {
		defer reader.Close()

		_, err = serviceClient.UploadStream(context.Background(), PERMANENT_CONTAINER_NAME, fmt.Sprintf("%d", id), reader, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for pager.More() {
		// advance to the next page
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return err
		}

		for _, blob := range page.Segment.BlobItems {
			// Downloads the chunk from the temporary container and writes it to the pipe
			blobDownloadResponse, err := serviceClient.DownloadStream(context.Background(), CONTAINER_NAME, *blob.Name, nil)
			if err != nil {
				log.Fatal(err)
			}

			bodyStream := blobDownloadResponse.Body
			_, err = io.Copy(writer, bodyStream)
			if err != nil {
				log.Fatal(err)
			}

			bodyStream.Close()
		}
	}

	writer.Close()

	return nil
}

func newUploadService() *uploadServiceImpl {
	return &uploadServiceImpl{}
}

func GetUploadService() UploadService {
	initUploadServiceOnce.Do(func() {
		uploadService = newUploadService()
	})

	return uploadService
}
