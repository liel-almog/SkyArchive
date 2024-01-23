package middleware

import (
	"errors"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lielalmog/file-uploader/backend/configs"
)

// Protected protect routes
func Protected() fiber.Handler {
	secret, err := configs.GetEnv("JWT_SECRET")
	if err != nil {
		panic(err)
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing or malformed JWT")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}
