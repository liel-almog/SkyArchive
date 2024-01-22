package controllers

import (
	"errors"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
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

	token, err := a.authService.Signup(signup)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fiber.ErrConflict
			}
		}

		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (a *authControllerImpl) Login(c *fiber.Ctx) error {
	login := new(models.Login)

	if err := c.BodyParser(login); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(login); err != nil {
		return fiber.ErrBadRequest
	}

	token, err := a.authService.Login(login)
	if err != nil {
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
