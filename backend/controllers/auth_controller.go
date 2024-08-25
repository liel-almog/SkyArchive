package controllers

import (
	"errors"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/liel-almog/SkyArchive/backend/configs"
	"github.com/liel-almog/SkyArchive/backend/errors/apperrors"
	"github.com/liel-almog/SkyArchive/backend/models"
	"github.com/liel-almog/SkyArchive/backend/services"
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
	signup := new(models.AuthSignup)

	if err := c.BodyParser(signup); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(signup); err != nil {
		return fiber.ErrBadRequest
	}

	token, err := a.authService.Signup(c.Context(), signup)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "User already exists")
		}

		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (a *authControllerImpl) Login(c *fiber.Ctx) error {
	login := new(models.AuthLogin)

	if err := c.BodyParser(login); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(login); err != nil {
		return fiber.ErrBadRequest
	}

	token, err := a.authService.Login(c.Context(), login)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		if errors.Is(err, apperrors.ErrUserNotFound) {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
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
