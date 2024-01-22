package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/lielalmog/file-uploader/backend/kafka"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
	segKafka "github.com/segmentio/kafka-go"
)

type UploadService interface {
	StartUpload(fileMetadate *models.FileMetadateDTO) (*int64, error)
	UploadChunk(fileHeader *multipart.FileHeader, id int64, chunkIndex int) error
	CompleteUploadEvent(id int64) error
}

type uploadServiceImpl struct {
	uploadRepository repositories.UploadRepository
}

var (
	initUploadServiceOnce sync.Once
	uploadService         *uploadServiceImpl
)

const (
	tempContainerName      = "temp-files"
	permanentContainerName = "permanent-files"
)

func (u *uploadServiceImpl) StartUpload(fileMetadate *models.FileMetadateDTO) (*int64, error) {
	return u.uploadRepository.SaveFileMetadata(fileMetadate)
}

func (u *uploadServiceImpl) UploadChunk(fileHeader *multipart.FileHeader, id int64, chunkIndex int) error {
	connectionString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	blobName := fmt.Sprintf("%d/%d", id, chunkIndex)
	serviceClient.UploadStream(context.Background(), tempContainerName, blobName, file, nil)
	if err != nil {
		return err
	}

	return nil
}

func (u *uploadServiceImpl) CompleteUploadEvent(id int64) error {
	writer := kafka.GetKafkaProducer()

	err := writer.WriteMessages(context.Background(), segKafka.Message{
		Value: []byte(fmt.Sprintf("%d", id)),
		Topic: kafka.UploadFilesTopic,
	})

	return err
}

func newUploadService() *uploadServiceImpl {
	return &uploadServiceImpl{
		uploadRepository: repositories.GetUploadRepository(),
	}
}

func GetUploadService() UploadService {
	initUploadServiceOnce.Do(func() {
		uploadService = newUploadService()
	})

	return uploadService
}
