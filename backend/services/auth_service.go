package services

import (
	"context"
	"sync"

	"github.com/lielalmog/SkyArchive/backend/configs"
	"github.com/lielalmog/SkyArchive/backend/errors/apperrors"
	"github.com/lielalmog/SkyArchive/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(ctx context.Context, signup *models.AuthSignup) (*string, error)
	Login(ctx context.Context, login *models.AuthLogin) (*string, error)
}

type authServiceImpl struct {
	userService UserService
	jwtService  JWTService
}

var (
	initAuthServiceOnce sync.Once
	authService         *authServiceImpl
)

func (a *authServiceImpl) Signup(ctx context.Context, signup *models.AuthSignup) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 14)
	if err != nil {
		return nil, err
	}

	signup.Password = string(bytes)

	id, err := a.userService.SaveUser(ctx, signup)
	if err != nil {
		return nil, err
	}

	claims := configs.CustomJwtClaims{
		Email:    signup.Email,
		Username: signup.Username,
		Id:       *id,
	}

	token, err := a.jwtService.GenerateToken(&claims)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *authServiceImpl) Login(ctx context.Context, login *models.AuthLogin) (*string, error) {
	user, err := a.userService.GetUserByEmail(ctx, &login.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	claims := configs.CustomJwtClaims{
		Email:    user.Email,
		Username: user.Username,
		Id:       user.ID,
	}

	token, err := a.jwtService.GenerateToken(&claims)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func newAuthService() *authServiceImpl {
	return &authServiceImpl{
		userService: GetUserService(),
		jwtService:  GetJWTService(),
	}
}

func GetAuthService() AuthService {
	initAuthServiceOnce.Do(func() {
		authService = newAuthService()
	})

	return authService
}
