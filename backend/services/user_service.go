package services

import (
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lielalmog/file-uploader/backend/errors/apperrors"
	"github.com/lielalmog/file-uploader/backend/errors/pgerrors"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/repositories"
)

type UserService interface {
	// GetUserByEmail returns a user by email. If no user is found, nil is returned without error.
	//
	// If an error occurs, the error is returned.
	GetUserByEmail(email *string) (*models.User, error)

	SaveUser(user *models.AuthSignup) (*int64, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

var (
	initUserServiceOnce sync.Once
	userService         *userServiceImpl
)

func (service *userServiceImpl) GetUserByEmail(email *string) (*models.User, error) {
	user, err := service.userRepository.FindUserByEmail(email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *userServiceImpl) SaveUser(user *models.AuthSignup) (*int64, error) {
	id, err := service.userRepository.SaveUser(user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrors.UniqueViolation {
				return nil, apperrors.ErrUserAlreadyExists
			}
		}

		return nil, err
	}

	return id, nil
}

func newUserService() *userServiceImpl {
	return &userServiceImpl{
		userRepository: repositories.GetUserRepository(),
	}
}

func GetUserService() UserService {
	initUserServiceOnce.Do(func() {
		userService = newUserService()
	})

	return userService
}
