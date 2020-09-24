package usecase

import (
	"context"
	"errors"
)

// All good validator mock
type trueValidator struct{}

func (t *trueValidator) ValidateEmail(email string) error {
	return nil
}
func (t *trueValidator) ValidatePassword(password string) error {
	return nil
}

// Invalid email mock
type falseValidatorEmail struct{}

func (f *falseValidatorEmail) ValidateEmail(email string) error {
	return &InvalidEmailError{errors.New("Bad"), email}
}
func (f *falseValidatorEmail) ValidatePassword(password string) error {
	return nil
}

// Invalid password mock
type falseValidatorPassword struct{}

func (f *falseValidatorPassword) ValidateEmail(email string) error {
	return nil
}
func (f *falseValidatorPassword) ValidatePassword(password string) error {
	return InvalidPasswordError
}

// UserStore lame mock
type happyPathUserStoreMock struct{}

func (m *happyPathUserStoreMock) Create(ctx context.Context, u *User) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) Update(ctx context.Context, u *User) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) ReadOne(ctx context.Context, crit ...Lookup) (*UserDto, error) {
	return &UserDto{
		Email:     "mock@mock.com",
		Username:  "mocking",
		FirstName: "Mocker",
		LastName:  "Mogul",
	}, nil
}

// Errored userStore
type erroredUserStoreMock struct{}

func (m *erroredUserStoreMock) Create(ctx context.Context, u *User) (*UserDto, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) Update(ctx context.Context, u *User) (*UserDto, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) ReadOne(ctx context.Context, crit ...Lookup) (*UserDto, error) {
	return nil, errors.New("Not found")
}

// Not found for ReadOne mock

type notFoundForReadOneStoreMock struct{}

func (m *notFoundForReadOneStoreMock) Create(ctx context.Context, u *User) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil

}

func (m *notFoundForReadOneStoreMock) Update(ctx context.Context, u *User) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *notFoundForReadOneStoreMock) ReadOne(ctx context.Context, c ...Lookup) (*UserDto, error) {
	return nil, errors.New("Not found")
}
