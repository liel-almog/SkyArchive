package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/controllers"
)

func NewFileRouter(router fiber.Router) {
	group := router.Group("/file")
	contoller := controllers.GetFileController()

	group.Get("/", contoller.GetUserFiles)
	group.Patch("/:id/favorite", contoller.UpdateFavorite)
	group.Patch("/:id/display-name", contoller.UpdateDisplayName)
	group.Delete("/:id", contoller.DeleteFile)
	group.Get("/download/:id", contoller.DownloadFile)

	uploadGroup := group.Group("/upload")

	uploadGroup.Post("/start", contoller.StartFileUpload)
	uploadGroup.Post("/complete/:id", contoller.CompleteFileUpload)
}
