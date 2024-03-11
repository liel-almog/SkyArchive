package models

import "time"

type File struct {
	FileID       int64     `json:"fileId" db:"file_id"`
	UserID       int64     `json:"userId" db:"user_id"`
	DisplayName  string    `json:"displayName" db:"display_name"`
	OriginalName string    `json:"originalName" db:"original_name"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mimeType" db:"mime_type"`
	Favorite     bool      `json:"favorite"`
	Status       string    `json:"status"`
	UploadedAt   time.Time `json:"uploadedAt" db:"uploaded_at"`
}

type UploadFileMetadateDTO struct {
	FileName string `json:"fileName" binding:"required" validate:"min=1"`
	Size     int64  `json:"size" binding:"required" validate:"min=1"`
	MimeType string `json:"mimeType" binding:"required" validate:"min=1"`
}

type FileMetadata struct {
	// embed the UploadFileMetadateDTO
	UploadFileMetadateDTO
	UserID int64 `json:"userId" binding:"required" validate:"min=1"`
}

type FileResDTO struct {
	FileID      int64     `json:"fileId"`
	DisplayName string    `json:"displayName"`
	UploadedAt  time.Time `json:"uploadedAt" db:"uploaded_at"`
	Favorite    bool      `json:"favorite"`
	Size        int64     `json:"size"`
	Status      string    `json:"status"`
}

type UpdateFavoriteDTO struct {
	Favorite bool `json:"favorite" binding:"required" validate:"boolean"`
}

type UpdateDisplayNameDTO struct {
	DisplayName string `json:"displayName" binding:"required" validate:"min=1"`
}
