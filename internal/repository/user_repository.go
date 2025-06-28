package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction on create user: %w", err)
	}
	defer tx.Rollback()

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (id, first_name, last_name, middle_na,e, date_of_birth, mobile_number, gender, email, password_hash, status, created_at)
    VALUES (:id, :first_name, :last_name, :middle_name, :date_of_birth, :mobile_number, :gender, :email, :password_hash, :status, :created_at)`

	_, err = tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (repo *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.QueryTimeout)
	defer cancel()

	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := repo.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
