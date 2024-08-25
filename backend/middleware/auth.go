package middleware

import (
	"encoding/json"
	"errors"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/liel-almog/SkyArchive/backend/configs"
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
		SuccessHandler: func(c *fiber.Ctx) error {
			mspClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
			customClaims := new(configs.CustomJwtClaims)

			jsonData, err := json.Marshal(mspClaims)
			if err != nil {
				return err
			}

			err = json.Unmarshal(jsonData, customClaims)
			if err != nil {
				return err
			}

			c.Locals("userClaims", customClaims)
			return c.Next()
		},
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
