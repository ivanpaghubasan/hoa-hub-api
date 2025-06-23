package repository

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}

func (repo *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
