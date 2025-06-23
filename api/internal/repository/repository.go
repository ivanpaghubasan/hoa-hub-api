package repository

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User interface {
		CreateUser(ctx context.Context, user *model.User) (*model.User, error)
		GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	}
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
