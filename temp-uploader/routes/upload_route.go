package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/be-file-streaming/controllers"
)

func NewUploadRouter(app fiber.Router) {
	group := app.Group("/upload")
	contoller := controllers.GetUploadController()

	group.Post("/chunk/start", contoller.StartUpload)
	group.Post("/chunk/complete/:id", contoller.CompleteUpload)
	group.Post("/chunk/:id/:chunkIndex", contoller.UploadChunk)
}
