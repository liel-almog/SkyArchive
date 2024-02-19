package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/controllers"
)

func NewAuthRouter(router fiber.Router) {
	group := router.Group("/auth")
	controller := controllers.GetAuthController()

	group.Post("/signup", controller.Signup)
	group.Post("/login", controller.Login)
}
