package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/lielalmog/file-uploader/backend/database"
	"github.com/lielalmog/file-uploader/backend/models"
)

type UploadRepository interface {
	SaveFileMetadata(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error)
}

type uploadRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initUploadRepositoryOnce sync.Once
	uploadRepository         *uploadRepositoryImpl
)

func (u *uploadRepositoryImpl) SaveFileMetadata(ctx context.Context, fileMetadate *models.UploadFileMetadateDTO) (*int64, error) {
	var id *int64

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := u.db.Pool.QueryRow(queryCtx,
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
