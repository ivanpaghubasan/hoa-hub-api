package repository

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
