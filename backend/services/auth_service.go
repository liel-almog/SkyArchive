package services

import (
	"sync"

	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(signup *models.Signup) error
}

type authServiceImpl struct {
	userRepository repositories.UserRepository
}

var (
	initAuthServiceOnce sync.Once
	authService         *authServiceImpl
)

func (a *authServiceImpl) Signup(signup *models.Signup) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 14)
	if err != nil {
		return err
	}

	signup.Password = string(bytes)

	if err := a.userRepository.SaveUser(signup); err != nil {
		return err
	}

	return nil
}

func newAuthService() *authServiceImpl {
	return &authServiceImpl{
		userRepository: repositories.GetUserRepository(),
	}
}

func GetAuthService() AuthService {
	initAuthServiceOnce.Do(func() {
		authService = newAuthService()
	})

	return authService
}
