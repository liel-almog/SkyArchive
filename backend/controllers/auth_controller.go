package controllers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/services"
)

type AuthController interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authControllerImpl struct {
	authService services.AuthService
}

var (
	initAuthControllerOnce sync.Once
	authController         *authControllerImpl
)

func (a *authControllerImpl) Signup(c *fiber.Ctx) error {
	signup := new(models.Signup)

	if err := c.BodyParser(signup); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(signup); err != nil {
		return fiber.ErrBadRequest
	}

	if err := a.authService.Signup(signup); err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (a *authControllerImpl) Login(c *fiber.Ctx) error {
	return nil
}

func newAuthController() *authControllerImpl {
	return &authControllerImpl{
		authService: services.GetAuthService(),
	}
}

func GetAuthController() AuthController {
	initAuthControllerOnce.Do(func() {
		authController = newAuthController()
	})

	return authController
}
