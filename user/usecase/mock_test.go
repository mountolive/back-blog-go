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

func (t *trueValidator) ValidatePasswordMatch(password, newPassword string) error {
	return nil
}

// Invalid email mock
type falseValidatorEmail struct{}

func (f *falseValidatorEmail) ValidateEmail(email string) error {
	return MalformedEmailError
}

func (f *falseValidatorEmail) ValidatePassword(password string) error {
	return nil
}

func (t *falseValidatorEmail) ValidatePasswordMatch(password, newPassword string) error {
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

func (t *falseValidatorPassword) ValidatePasswordMatch(password, newPassword string) error {
	return nil
}

// Passwords don't match mock
type falseValidatorPasswordsNotMatching struct{}

func (f *falseValidatorPasswordsNotMatching) ValidateEmail(email string) error {
	return nil
}

func (f *falseValidatorPasswordsNotMatching) ValidatePassword(password string) error {
	return nil
}

func (t *falseValidatorPasswordsNotMatching) ValidatePasswordMatch(password, newPassword string) error {
	return PasswordsDontMatchError
}

// UserStore lame mock
type happyPathUserStoreMock struct {
	email,
	username string
}

func (m *happyPathUserStoreMock) Create(ctx context.Context, u *CreateUserDto) (*User, error) {
	return &User{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) Update(ctx context.Context, i string, u *UpdateUserDto) (*User, error) {
	return &User{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*User, error) {
	var user *User
	if m.email != "" && m.username != "" {
		user = &User{
			Email:    m.email,
			Username: m.username,
		}
	}
	return user, nil
}

func (m *happyPathUserStoreMock) CheckIfCorrectPassword(ctx context.Context, d *CheckUserAndPasswordDto) error {
	return nil
}

func (m *happyPathUserStoreMock) UpdatePassword(c context.Context, d *ChangePasswordDto) error {
	return nil
}

// Errored userStore
type erroredUserStoreMock struct{}

func (m *erroredUserStoreMock) Create(ctx context.Context, u *CreateUserDto) (*User, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) Update(ctx context.Context, i string, u *UpdateUserDto) (*User, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*User, error) {
	return nil, errors.New("Something happened")
}

func (m *erroredUserStoreMock) CheckIfCorrectPassword(ctx context.Context, d *CheckUserAndPasswordDto) error {
	return UserPasswordNotMatchingError
}

func (m *erroredUserStoreMock) UpdatePassword(c context.Context, d *ChangePasswordDto) error {
	return nil
}
