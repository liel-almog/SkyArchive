package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lielalmog/SkyArchive/backend/database"
	"github.com/lielalmog/SkyArchive/backend/models"
)

type FileRepository interface {
	SaveFileMetadata(ctx context.Context, fileMetadate *models.FileMetadata) (*int64, error)
	GetUserFiles(ctx context.Context, userId *int64) ([]models.FileResDTO, error)
	UpdateFavorite(ctx context.Context, fileId *int64, userId *int64, updateFavorite *models.UpdateFavoriteDTO) error
	UpdateDisplayName(ctx context.Context, fileId *int64, userId *int64, updateDisplayName *models.UpdateDisplayNameDTO) error
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

func (u *fileRepositoryImpl) GetUserFiles(ctx context.Context, userId *int64) ([]models.FileResDTO, error) {
	var files []models.FileResDTO

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := u.db.Pool.Query(queryCtx,
		`SELECT file_id, display_name, uploaded_at, favorite, size, status
		 FROM files
		 WHERE user_id = $1
		 ORDER BY uploaded_at`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.FileResDTO])
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (u *fileRepositoryImpl) UpdateFavorite(ctx context.Context, fileId *int64, userId *int64, updateFavorite *models.UpdateFavoriteDTO) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := u.db.Pool.Exec(queryCtx,
		`UPDATE files 
		SET favorite = $1 
		WHERE file_id = $2 AND user_id = $3`,
		updateFavorite.Favorite, fileId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *fileRepositoryImpl) UpdateDisplayName(ctx context.Context, fileId *int64, userId *int64, updateDisplayName *models.UpdateDisplayNameDTO) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := u.db.Pool.Exec(queryCtx,
		`UPDATE files 
		SET display_name = $1 
		WHERE file_id = $2 AND user_id = $3`,
		updateDisplayName.DisplayName, fileId, userId)
	if err != nil {
		return err
	}

	return nil
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
