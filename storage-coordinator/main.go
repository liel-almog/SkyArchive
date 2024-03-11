package main

import (
	"sync"

	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
)

const (
	fileUplaodFinilizationTopic = "file-upload-finalization"
	fileDeleteTopic             = "file-delete"

	permanentContainerName = "permanent-files"
	backupContainerName    = "backup-files"
	tempContainerName      = "temp-files"
)

func main() {
	var wg sync.WaitGroup
	configs.InitEnv()

	connString, err := configs.GetEnv("AZURE_STORAGE_CONNECTION_STRING")
	if err != nil {
		panic(err)
	}

	wg.Add(2)

	go func() {
		fileUpload(connString)
		wg.Done()
	}()

	go func() {
		fileDelete(connString)
		wg.Done()
	}()

	wg.Wait()
}
