package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/configs"
	"github.com/lielalmog/SkyArchive/backend/kafka"
	"github.com/lielalmog/SkyArchive/backend/models"
	"github.com/lielalmog/SkyArchive/backend/repositories"
	segKafka "github.com/segmentio/kafka-go"
)

type FileService interface {
	SaveFileMetadata(ctx context.Context, fileMetadateDTO *models.FileMetadata) (*int64, error)
	GenerateSasToken(ctx context.Context, fileId *int64) (*string, error)
	CompleteFileUploadEvent(ctx context.Context, fileId *int64) error
	GetUserFiles(ctx context.Context, userId *int64) ([]models.FileResDTO, error)
	UpdateFavorite(ctx context.Context, fileId *int64, userId *int64, updateFavoriteDTO *models.UpdateFavoriteDTO) error
	UpdateDisplayName(ctx context.Context, fileId *int64, userId *int64, updateDisplayNameDTO *models.UpdateDisplayNameDTO) error
	DeleteFile(ctx context.Context, fileId *int64, userId *int64) error
	GetFileByUser(ctx context.Context, fileId *int64, userId *int64) (*models.File, error)
}

type fileServiceImpl struct {
	fileRepository repositories.FileRepository
}

var (
	initFileServiceOnce sync.Once
	fileService         *fileServiceImpl
)

const (
	tempContainerName      = "temp-files"
	permanentContainerName = "permanent-files"
)

func parseAzureStorageConnectionString(connectionString string) (accountName, accountKey string, err error) {
	parts := strings.Split(connectionString, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "AccountName=") {
			accountName = strings.TrimPrefix(part, "AccountName=")
		} else if strings.HasPrefix(part, "AccountKey=") {
			accountKey = strings.TrimPrefix(part, "AccountKey=")
		}
	}

	if accountName == "" || accountKey == "" {
		return "", "", errors.New("account name or key not found in connection string")
	}

	return accountName, accountKey, nil
}

func (u *fileServiceImpl) SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error) {
	return u.fileRepository.SaveFileMetadata(ctx, fileMetadate)
}

func (u *fileServiceImpl) GenerateSasToken(ctx context.Context, fileId *int64) (*string, error) {
	connectionString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	accountName, accountKey, err := parseAzureStorageConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	// extract the account name and key from the connection string
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	// create a SAS token that's valid for one hour.
	// since this is a network call, there's a timeout.
	sasPermissions := sas.BlobPermissions{Create: true, Add: true, Write: true, Read: true, Delete: true, List: true}
	fileName := fmt.Sprintf("%d", *fileId)

	token, err := sas.BlobSignatureValues{
		Protocol:    sas.ProtocolHTTPS,
		ExpiryTime:  time.Now().UTC().Add(3 * time.Hour),
		Permissions: sasPermissions.String(),

		// Start time 10 minutes ago to avoid clock skew.
		StartTime:     time.Now().UTC().Add(-10 * time.Minute),
		ContainerName: tempContainerName,
		BlobName:      fileName,
	}.SignWithSharedKey(credential)
	if err != nil {
		return nil, err
	}

	signedUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, tempContainerName, fileName, token.Encode())

	return &signedUrl, nil
}

func (u *fileServiceImpl) CompleteFileUploadEvent(ctx context.Context, fileId *int64) error {
	writer := kafka.GetKafkaProducer()

	payload, err := json.Marshal(models.KafkaFileUploadFinalizationMessage{
		FileId: fileId,
	})
	if err != nil {
		return err
	}

	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = writer.WriteMessages(writeCtx, segKafka.Message{
		Value: payload,
		Topic: kafka.FileUploadFinalizationTopic,
	})

	return err
}

func (u *fileServiceImpl) GetUserFiles(ctx context.Context, userId *int64) ([]models.FileResDTO, error) {
	return u.fileRepository.GetUserFiles(ctx, userId)
}

func (u *fileServiceImpl) UpdateFavorite(ctx context.Context, fileId *int64, userId *int64, updateFavoriteDTO *models.UpdateFavoriteDTO) error {
	n, err := u.fileRepository.UpdateFavorite(ctx, fileId, userId, updateFavoriteDTO)

	if err != nil {
		return err
	}

	if n != 1 {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}

	return nil
}

func (u *fileServiceImpl) UpdateDisplayName(ctx context.Context, fileId *int64, userId *int64, updateDisplayNameDTO *models.UpdateDisplayNameDTO) error {
	n, err := u.fileRepository.UpdateDisplayName(ctx, fileId, userId, updateDisplayNameDTO)

	if err != nil {
		return err
	}

	if n != 1 {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}

	return nil
}

func (u *fileServiceImpl) DeleteFile(ctx context.Context, fileId *int64, userId *int64) error {
	writer := kafka.GetKafkaProducer()

	kafkaResult := make(chan error)

	go func() {
		defer close(kafkaResult)

		payload, err := json.Marshal(models.KafkaFileUploadFinalizationMessage{
			FileId: fileId,
		})
		if err != nil {
			kafkaResult <- err
			return
		}

		writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = writer.WriteMessages(writeCtx, segKafka.Message{
			Value: payload,
			Topic: kafka.FileDeleteTopic,
		})

		if err != nil {
			kafkaResult <- err
			return
		}

		kafkaResult <- nil
	}()

	n, err := u.fileRepository.DeleteFile(ctx, kafkaResult, fileId, userId)
	if err != nil {
		return err
	}

	if n == 0 {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}

	return nil
}

func (u *fileServiceImpl) GetFileByUser(ctx context.Context, fileId *int64, userId *int64) (*models.File, error) {
	return u.fileRepository.GetFileByUser(ctx, fileId, userId)
}

func newFileService() *fileServiceImpl {
	return &fileServiceImpl{
		fileRepository: repositories.GetFileRepository(),
	}
}

func GetFileService() FileService {
	initFileServiceOnce.Do(func() {
		fileService = newFileService()
	})

	return fileService
}
