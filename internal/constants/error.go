package constants

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRecordExists    = errors.New("record exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInternalServer  = errors.New("internal server errror")
)
