package services

import (
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lielalmog/file-uploader/backend/configs"
)

type JWTService interface {
	GenerateToken(claims jwt.MapClaims) (*string, error)
}

type jwtServiceImpl struct{}

var (
	initJWTServiceOnce sync.Once
	jwtService         *jwtServiceImpl
)

func (j *jwtServiceImpl) GenerateToken(claims jwt.MapClaims) (*string, error) {
	//TODO Add default claims - exp, iat, iss, sub

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := configs.GetEnv("JWT_SECRET")
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func newJWTService() *jwtServiceImpl {
	return &jwtServiceImpl{}
}

func GetJWTService() JWTService {
	initJWTServiceOnce.Do(func() {
		jwtService = newJWTService()
	})

	return jwtService
}
