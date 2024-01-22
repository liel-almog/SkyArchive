package repositories

import (
	"context"
	"sync"

	"github.com/lielalmog/file-uploader/backend/database"
	"github.com/lielalmog/file-uploader/backend/models"
)

type UploadRepository interface {
	SaveFileMetadata(fileMetadate *models.UploadFileMetadateDTO) (*int64, error)
}

type uploadRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initUploadRepositoryOnce sync.Once
	uploadRepository         *uploadRepositoryImpl
)

func (u *uploadRepositoryImpl) SaveFileMetadata(fileMetadate *models.UploadFileMetadateDTO) (*int64, error) {
	var id *int64

	row := u.db.Pool.QueryRow(context.Background(),
		"INSERT INTO files (display_name, original_name, size) VALUES ($1, $1, $2) RETURNING file_id", fileMetadate.FileName, fileMetadate.Size)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return id, nil
}

func newUploadRepository() *uploadRepositoryImpl {
	return &uploadRepositoryImpl{
		db: database.GetDB(),
	}
}

func GetUploadRepository() UploadRepository {
	initUploadRepositoryOnce.Do(func() {
		uploadRepository = newUploadRepository()
	})

	return uploadRepository
}
