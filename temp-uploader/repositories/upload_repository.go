package repositories

import (
	"context"
	"sync"

	"github.com/lielalmog/be-file-streaming/database"
	"github.com/lielalmog/be-file-streaming/models"
)

type UploadRepository interface {
	SaveFileMetadata(fileMetadate *models.FileMetadateDTO) (*int64, error)
}

type UploadRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initUploadRepositoryOnce sync.Once
	uploadRepository         *UploadRepositoryImpl
)

func (u *UploadRepositoryImpl) SaveFileMetadata(fileMetadate *models.FileMetadateDTO) (*int64, error) {
	var id *int64

	row := u.db.Pool.QueryRow(context.Background(),
		"INSERT INTO files (display_name, original_name, size) VALUES ($1, $1, $2) RETURNING file_id", fileMetadate.FileName, fileMetadate.Size)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return id, nil
}

func newUploadRepository() *UploadRepositoryImpl {
	return &UploadRepositoryImpl{
		db: database.GetDB(),
	}
}

func GetUploadRepository() UploadRepository {
	initUploadRepositoryOnce.Do(func() {
		uploadRepository = newUploadRepository()
	})

	return uploadRepository
}
