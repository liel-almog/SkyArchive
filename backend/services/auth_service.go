package services

import (
	"sync"

	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(signup *models.Signup) (*string, error)
	Login(login *models.Login) (*string, error)
}

type authServiceImpl struct {
	userRepository repositories.UserRepository
	jwtService     JWTService
}

var (
	initAuthServiceOnce sync.Once
	authService         *authServiceImpl
)

func (a *authServiceImpl) Signup(signup *models.Signup) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 14)
	if err != nil {
		return nil, err
	}

	signup.Password = string(bytes)

	if err := a.userRepository.SaveUser(signup); err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(map[string]interface{}{
		"email":    signup.Email,
		"username": signup.Username,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *authServiceImpl) Login(login *models.Login) (*string, error) {
	user, err := a.userRepository.FindUserByEmail(login.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func newAuthService() *authServiceImpl {
	return &authServiceImpl{
		userRepository: repositories.GetUserRepository(),
		jwtService:     GetJWTService(),
	}
}

func GetAuthService() AuthService {
	initAuthServiceOnce.Do(func() {
		authService = newAuthService()
	})

	return authService
}
