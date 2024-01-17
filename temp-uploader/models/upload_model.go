package models

type FileMetadateDTO struct {
	FileName string `json:"fileName" binding:"required" validate:"min=1"`
	Size     int64  `json:"size" binding:"required" validate:"min=1"`
}
