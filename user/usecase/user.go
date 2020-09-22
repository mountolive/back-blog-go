package usecase

import (
	"errors"
	"fmt"
)

var InvalidPasswordError = errors.New("Password should have " +
	"more than 6 characters, " +
	"at least one number and " +
	"one letter")

type UserValidator interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
}

type User struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type InvalidEmailError struct {
	Err      error
	BadEmail string
}

func (u *InvalidEmailError) Error() string {
	return fmt.Sprintf("The passed email: %v, has an error: %v \n",
		u.BadEmail, u.Err)
}

func (u *InvalidEmailError) Unwrap() error {
	return u.Err
}

func NewUser(username, email, password string, validator UserValidator) (*User, error) {
	user := &User{}
	user.Username = username
	err := validator.ValidateEmail(email)
	if err != nil {
		return user, fmt.Errorf("The email is invalid. %w", err)
	}
	user.Email = email
	err = validator.ValidatePassword(password)
	if err != nil {
		return user, fmt.Errorf("The password is invalid. %w", err)
	}
	user.Password = password
	return user, nil
}
