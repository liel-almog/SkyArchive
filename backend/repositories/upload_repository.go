package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/lielalmog/SkyArchive/backend/database"
	"github.com/lielalmog/SkyArchive/backend/models"
)

type UploadRepository interface {
	SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error)
}

type uploadRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initUploadRepositoryOnce sync.Once
	uploadRepository         *uploadRepositoryImpl
)

func (u *uploadRepositoryImpl) SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error) {
	var id *int64

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := u.db.Pool.QueryRow(queryCtx,
		`INSERT INTO files (display_name, original_name, size, user_id, mime_type) 
		VALUES ($1, $1, $2, $3, $4) RETURNING file_id`,
		fileMetadate.FileName, fileMetadate.Size, fileMetadate.UserID, fileMetadate.MimeType)

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
