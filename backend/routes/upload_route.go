package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/file-uploader/backend/controllers"
)

func NewUploadRouter(router fiber.Router) {
	group := router.Group("/upload")
	contoller := controllers.GetUploadController()

	group.Post("/chunk/start", contoller.StartUpload)
	group.Post("/chunk/complete/:id", contoller.CompleteUpload)
	group.Post("/chunk/:id/:chunkIndex", contoller.UploadChunk)
}
