package user

import (
	"errors"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrUserNotFound     = errors.New("user not found")
)
