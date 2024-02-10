package services

import (
	"encoding/json"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lielalmog/file-uploader/backend/configs"
)

type JWTService interface {
	GenerateToken(claims *configs.CustomJwtClaims) (*string, error)
	ExtractClaims(token *jwt.Token) (*configs.CustomJwtClaims, error)
}

type jwtServiceImpl struct{}

var (
	initJWTServiceOnce sync.Once
	jwtService         *jwtServiceImpl
)

func (j *jwtServiceImpl) GenerateToken(claims *configs.CustomJwtClaims) (*string, error) {
	//TODO Add default claims - exp, iat, iss, sub
	b, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	var claimsMap jwt.MapClaims
	err = json.Unmarshal(b, &claimsMap)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)

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

func (j *jwtServiceImpl) ExtractClaims(token *jwt.Token) (*configs.CustomJwtClaims, error) {
	mspClaims := token.Claims.(jwt.MapClaims)
	customClaims := new(configs.CustomJwtClaims)

	jsonData, err := json.Marshal(mspClaims)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, customClaims)
	if err != nil {
		return nil, err
	}

	return customClaims, nil
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
