package domain

import (
	"errors"
	"net/mail"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidEmail       = errors.New("email is invalid")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewUser(
	id,
	email,
	password string,
) (*User, error) {

	user := &User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}

	return nil
}
