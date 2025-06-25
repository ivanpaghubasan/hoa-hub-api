package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
	"github.com/jmoiron/sqlx"
)

const (
	QueryTimeout = time.Second * 5
	DateFormat   = "2006-01-02"

	ActiveStatus   = "active"
	InactiveStatus = "inactive"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrRecordExists   = errors.New("record exists")
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
