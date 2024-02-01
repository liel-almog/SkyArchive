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
	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/lielalmog/file-uploader/backend/kafka"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
	segKafka "github.com/segmentio/kafka-go"
)

type UploadService interface {
	SaveFileMetadata(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error)
	GenerateSasToken(ctx context.Context, fileId *int64) (*string, error)
	CompleteUploadEvent(ctx context.Context, fileId *int64) error
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

func (u *uploadServiceImpl) SaveFileMetadata(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error) {
	return u.uploadRepository.SaveFileMetadata(ctx, fileMetadate)
}

func (u *uploadServiceImpl) GenerateSasToken(ctx context.Context, fileId *int64) (*string, error) {
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

func (u *uploadServiceImpl) CompleteUploadEvent(ctx context.Context, fileId *int64) error {
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
