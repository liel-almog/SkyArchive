package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/controllers"
)

func NewFileRouter(router fiber.Router) {
	group := router.Group("/file")
	contoller := controllers.GetFileController()

	group.Get("/", contoller.GetUserFiles)

	uploadGroup := group.Group("/upload")

	uploadGroup.Post("/start", contoller.StartFileUpload)
	uploadGroup.Post("/complete/:id", contoller.CompleteFileUpload)
}
