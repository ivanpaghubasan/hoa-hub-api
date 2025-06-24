package repository

import (
	"context"
	"time"

	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/jmoiron/sqlx"
)

const (
	QueryTimeout = time.Second * 5
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type Repository struct {
	UserRepository UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
