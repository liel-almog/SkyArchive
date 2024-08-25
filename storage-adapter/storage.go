package main

import "context"

// StorageProvider is an interface for storage operations
type StorageProvider interface {
	UploadFile(ctx context.Context, fileName string, data []byte) (string, error)
	DeleteFile(ctx context.Context, fileName string) error
	GenerateSasToken(ctx context.Context, fileId *int64) (*string, error)
}

type storageProviderImpl struct {
	// implementation
}

func (s *storageProviderImpl) UploadFile(ctx context.Context, fileName string, data []byte) (string, error) {
	// implementation

	return "", nil
}
