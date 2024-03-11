package main

type FileUploadFinalizationMessage struct {
	FileId *int64 `json:"fileId"`
}

type FileDeleteMessage struct {
	FileId *int64 `json:"fileId"`
}
