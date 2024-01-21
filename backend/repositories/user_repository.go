package repositories

import (
	"context"
	"sync"

	"github.com/lielalmog/file-uploader/backend/database"
	"github.com/lielalmog/file-uploader/backend/models"
)

type UserRepository interface {
	SaveUser(*models.Signup) error
	FindUserByEmail() error
}

type userRepositoryImpl struct {
	db *database.PostgreSQLpgx
}

var (
	initUserRepositoryOnce sync.Once
	userRepository         *userRepositoryImpl
)

func (repo *userRepositoryImpl) SaveUser(signup *models.Signup) error {
	_, err := repo.db.Pool.Exec(context.Background(), "INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", signup.Email, signup.Password, signup.Username)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepositoryImpl) FindUserByEmail() error {
	return nil
}

func newUserRepository() *userRepositoryImpl {
	return &userRepositoryImpl{
		db: database.GetDB(),
	}
}

func GetUserRepository() UserRepository {
	initUserRepositoryOnce.Do(func() {
		userRepository = newUserRepository()
	})

	return userRepository
}
