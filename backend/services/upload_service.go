package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/lielalmog/file-uploader/backend/kafka"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
	segKafka "github.com/segmentio/kafka-go"
)

type UploadService interface {
	StartUpload(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error)
	UploadChunk(ctx context.Context, fileHeader *multipart.FileHeader, id int64, chunkIndex int) error
	CompleteUploadEvent(ctx context.Context, id int64) error
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

func (u *uploadServiceImpl) StartUpload(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error) {
	return u.uploadRepository.SaveFileMetadata(ctx, fileMetadate)
}

func (u *uploadServiceImpl) UploadChunk(ctx context.Context, fileHeader *multipart.FileHeader, id int64, chunkIndex int) error {
	connectionString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	uploadCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	blobName := fmt.Sprintf("%d/%d", id, chunkIndex)
	serviceClient.UploadStream(uploadCtx, tempContainerName, blobName, file, nil)
	if err != nil {
		return err
	}

	return nil
}

func (u *uploadServiceImpl) CompleteUploadEvent(ctx context.Context, id int64) error {
	writer := kafka.GetKafkaProducer()

	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := writer.WriteMessages(writeCtx, segKafka.Message{
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
