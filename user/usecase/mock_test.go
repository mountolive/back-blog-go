package usecase

import (
	"context"
	"errors"
	"fmt"
)

// Mock logger
type mockLogger struct{}

func (m *mockLogger) LogError(err error) {
	fmt.Printf("%v \n", err)
}

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

func (m *happyPathUserStoreMock) Create(ctx context.Context, u *CreateUserDto) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) Update(ctx context.Context, u *UpdateUserDto) (*UserDto, error) {
	return &UserDto{
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (m *happyPathUserStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*UserDto, error) {
	var user *UserDto
	if m.email != "" && m.username != "" {
		user = &UserDto{
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

func (m *erroredUserStoreMock) Create(ctx context.Context, u *CreateUserDto) (*UserDto, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) Update(ctx context.Context, u *UpdateUserDto) (*UserDto, error) {
	return nil, errors.New("Not found")
}

func (m *erroredUserStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*UserDto, error) {
	return nil, UserNotFoundError
}

func (m *erroredUserStoreMock) CheckIfCorrectPassword(ctx context.Context, d *CheckUserAndPasswordDto) error {
	return UserPasswordNotMatchingError
}

func (m *erroredUserStoreMock) UpdatePassword(c context.Context, d *ChangePasswordDto) error {
	return nil
}

// Incorrect match of old password mock

type oldPasswordNotMatchingStoreMock struct{}

func (m *oldPasswordNotMatchingStoreMock) Create(ctx context.Context, u *CreateUserDto) (*UserDto, error) {

	return nil, nil
}

func (m *oldPasswordNotMatchingStoreMock) Update(ctx context.Context, u *UpdateUserDto) (*UserDto, error) {
	return nil, nil
}

func (m *oldPasswordNotMatchingStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*UserDto, error) {
	invalid := "invalid"
	return &UserDto{
		Email:     invalid,
		Username:  invalid,
		FirstName: invalid,
		LastName:  invalid,
	}, nil

}

func (m *oldPasswordNotMatchingStoreMock) CheckIfCorrectPassword(ctx context.Context, d *CheckUserAndPasswordDto) error {
	return nil
}

func (m *oldPasswordNotMatchingStoreMock) UpdatePassword(c context.Context, d *ChangePasswordDto) error {
	return nil
}

// Wrong data on found for ReadOne

type incorrectFoundForReadOneStoreMock struct{}

func (m *incorrectFoundForReadOneStoreMock) Create(ctx context.Context, u *CreateUserDto) (*UserDto, error) {
	return nil, nil
}

func (m *incorrectFoundForReadOneStoreMock) Update(ctx context.Context, u *UpdateUserDto) (*UserDto, error) {
	return nil, nil
}

func (m *incorrectFoundForReadOneStoreMock) ReadOne(ctx context.Context, filter *ByUsernameOrEmail) (*UserDto, error) {
	invalid := "invalid"
	return &UserDto{
		Email:     invalid,
		Username:  invalid,
		FirstName: invalid,
		LastName:  invalid,
	}, nil

}

func (m *incorrectFoundForReadOneStoreMock) CheckIfCorrectPassword(ctx context.Context, d *CheckUserAndPasswordDto) error {
	return nil
}

func (m *incorrectFoundForReadOneStoreMock) UpdatePassword(c context.Context, d *ChangePasswordDto) error {
	return nil
}
