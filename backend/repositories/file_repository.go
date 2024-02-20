package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/lielalmog/SkyArchive/backend/database"
	"github.com/lielalmog/SkyArchive/backend/models"
)

type FileRepository interface {
	SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error)
}

type fileRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initFileRepositoryOnce sync.Once
	fileRepository         *fileRepositoryImpl
)

func (u *fileRepositoryImpl) SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error) {
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

func newFileRepository() *fileRepositoryImpl {
	return &fileRepositoryImpl{
		db: database.GetDB(),
	}
}

func GetFileRepository() FileRepository {
	initFileRepositoryOnce.Do(func() {
		fileRepository = newFileRepository()
	})

	return fileRepository
}
