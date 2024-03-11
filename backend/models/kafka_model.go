package models

type KafkaFileUploadFinalizationMessage struct {
	FileId *int64 `json:"fileId"`
}

type KafkaFileDeleteMessage struct {
	FileId *int64 `json:"fileId"`
}
