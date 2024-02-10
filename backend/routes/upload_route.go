package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/controllers"
)

func NewUploadRouter(router fiber.Router) {
	group := router.Group("/upload")
	contoller := controllers.GetUploadController()

	group.Post("/start", contoller.StartUpload)
	group.Post("/complete/:id", contoller.CompleteUpload)
}
