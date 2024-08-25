package services

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/liel-almog/SkyArchive/backend/errors/apperrors"
	"github.com/liel-almog/SkyArchive/backend/errors/pgerrors"
	"github.com/liel-almog/SkyArchive/backend/models"
	"github.com/liel-almog/SkyArchive/backend/repositories"
)

type UserService interface {
	// GetUserByEmail returns a user by email. If no user is found, nil is returned without error.
	//
	// If an error occurs, the error is returned.
	GetUserByEmail(ctx context.Context, email *string) (*models.User, error)

	SaveUser(ctx context.Context, user *models.AuthSignup) (*int64, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

var (
	initUserServiceOnce sync.Once
	userService         *userServiceImpl
)

func (service *userServiceImpl) GetUserByEmail(ctx context.Context, email *string) (*models.User, error) {
	user, err := service.userRepository.FindUserByEmail(ctx, email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *userServiceImpl) SaveUser(ctx context.Context, user *models.AuthSignup) (*int64, error) {
	id, err := service.userRepository.SaveUser(ctx, user)
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
