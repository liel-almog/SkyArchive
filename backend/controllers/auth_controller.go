package controllers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authControllerImpl struct{}

var (
	initAuthControllerOnce sync.Once
	authController         *authControllerImpl
)

func (a *authControllerImpl) Signup(c *fiber.Ctx) error {
	return nil
}

func (a *authControllerImpl) Login(c *fiber.Ctx) error {
	return nil
}

func newAuthController() *authControllerImpl {
	return &authControllerImpl{}
}

func GetAuthController() AuthController {
	initAuthControllerOnce.Do(func() {
		authController = newAuthController()
	})

	return authController
}
