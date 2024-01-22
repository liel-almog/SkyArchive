package repositories

import (
	"context"
	"sync"

	"github.com/lielalmog/file-uploader/backend/database"
	"github.com/lielalmog/file-uploader/backend/models"
)

type UserRepository interface {
	SaveUser(*models.Signup) error
	FindUserByEmail(email string) (*models.User, error)
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

func (repo *userRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	user := new(models.User)

	err := repo.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, username FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Username)

	if err != nil {
		return nil, err
	}

	return user, nil
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
