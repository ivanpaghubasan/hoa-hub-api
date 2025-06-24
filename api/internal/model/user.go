package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `db:"id"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	MiddleName   *string    `db:"middle_name"`
	DateOfBirth  *time.Time `db:"date_of_birth"`
	Gender       string     `db:"gender"`
	MobileNumber string     `db:"mobile_number"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type Account struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
