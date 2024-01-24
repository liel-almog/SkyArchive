package services

import (
	"sync"

	"github.com/lielalmog/file-uploader/backend/errors/apperrors"
	"github.com/lielalmog/file-uploader/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(signup *models.AuthSignup) (*string, error)
	Login(login *models.AuthLogin) (*string, error)
}

type authServiceImpl struct {
	userService UserService
	jwtService  JWTService
}

var (
	initAuthServiceOnce sync.Once
	authService         *authServiceImpl
)

func (a *authServiceImpl) Signup(signup *models.AuthSignup) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 14)
	if err != nil {
		return nil, err
	}

	signup.Password = string(bytes)

	id, err := a.userService.SaveUser(signup)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(map[string]interface{}{
		"email":    signup.Email,
		"username": signup.Username,
		"id":       id,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *authServiceImpl) Login(login *models.AuthLogin) (*string, error) {
	user, err := a.userService.GetUserByEmail(&login.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := a.jwtService.GenerateToken(map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
		"id":       user.ID,
	})
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
