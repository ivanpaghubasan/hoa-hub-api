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
	MobileNumber string     `db:"mobile_number"`
	Gender       string     `db:"gender"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
