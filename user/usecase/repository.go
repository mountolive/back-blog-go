// This package defines all the regular use cases related to users
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id        string
	Email     string
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserDto struct {
	Email,
	Password,
	RepeatedPassword,
	Username,
	FirstName,
	LastName string
}

type UpdateUserDto struct {
	Email,
	Username,
	FirstName,
	LastName string
}

type ChangePasswordDto struct {
	Email,
	Username,
	OldPassword,
	NewPassword,
	RepeatedPassword string
}

type CheckUserAndPasswordDto struct {
	Email,
	Username,
	Password string
}

type ByUsernameOrEmail struct {
	Username,
	Email string
}

// Common errors
var (
	EmailOrUsernameAlreadyInUseError = errors.New(
		"The email or username passed are already in use")
	UserNotFoundError = errors.New(
		"The user was not found in the DB")
	MalformedEmailError = errors.New(
		"The email passed is invalid")
	UserPasswordNotMatchingError = errors.New(
		"Seems like user/password data doesn't match with DB record")
	PasswordsDontMatchError = errors.New(
		"Password and RepeatedPassword don't match")
	InvalidPasswordError = errors.New(
		"Password doesn't comply with expected structure")
	OperationCanceledError = errors.New(
		"The context of the operation was canceled")
	CorruptedStoreError = errors.New(
		"The UserStore used is returning inconsistent results")
)

const unknownErrorInStore = "Found reported from store: %s and %s, but wrong dto returned"

// Contract for the needs of the repository in terms of persistance:
//    Defines which methods would be needed for each usecase
type UserStore interface {
	Create(context.Context, *CreateUserDto) (*User, error)
	Update(context.Context, string, *UpdateUserDto) (*User, error)
	UpdatePassword(context.Context, *ChangePasswordDto) error
	ReadOne(context.Context, *ByUsernameOrEmail) (*User, error)
	CheckIfCorrectPassword(context.Context, *CheckUserAndPasswordDto) error
}

// Validator for email's strings
type EmailValidator interface {
	ValidateEmail(email string) error
}

// Validator for password's constrains
type PasswordValidator interface {
	ValidatePassword(password string) error
}

// Validates equality of passwords
type MatchValidator interface {
	ValidatePasswordMatch(password, repeatedPassword string) error
}

// Contract for the needs of the repository in terms of validation:
//     Methods needed by each usecase for validating the user's data
type UserValidator interface {
	EmailValidator
	PasswordValidator
	MatchValidator
}

// Repository defines the basic usecases for the users' domain
type Repository interface {
	ReadUser(ctx context.Context, loginCred string) (*User, error)
	ChangePassword(ctx context.Context, changePass *ChangePasswordDto) error
	CreateUser(ctx context.Context, user *CreateUserDto) (*User, error)
	UpdateUser(ctx context.Context, id string, user *UpdateUserDto) (*User, error)
}

// Basic repository struct. Store is used for persitance and Validator
// for field validation
type UserRepository struct {
	Store     UserStore
	Validator UserValidator
}

var _ Repository = &UserRepository{}

// Reads an user either by her Username or by her Email
func (r *UserRepository) ReadUser(
	ctx context.Context,
	loginCred string,
) (*User, error) {
	byUsernameOrEmail := &ByUsernameOrEmail{
		Username: loginCred,
		Email:    loginCred,
	}
	return r.Store.ReadOne(ctx, byUsernameOrEmail)
}

// Changes password and persists. Returns an error on validation or
// store's retrieval/persistence
func (r *UserRepository) ChangePassword(
	ctx context.Context,
	changePass *ChangePasswordDto,
) error {
	err := r.Validator.ValidateEmail(changePass.Email)
	if err != nil {
		return logErrorAndWrap(err,
			"An error occurred when validating the user's email, ChangePassword")
	}
	toCheck := &CheckUserAndPasswordDto{
		Email:    changePass.Email,
		Username: changePass.Username,
		Password: changePass.OldPassword,
	}
	err = r.Store.CheckIfCorrectPassword(ctx, toCheck)
	if err != nil {
		return logErrorAndWrap(err, "An error occurred on the UserStore, ChangePassword")
	}
	err = r.validatePasswords(changePass.NewPassword, changePass.RepeatedPassword)
	if err != nil {
		// This error is already wrapped by the validatePasswords function
		return err
	}
	return r.Store.UpdatePassword(ctx, changePass)
}

// Creates an user. Returns an error on validation
func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *CreateUserDto,
) (*User, error) {
	err := r.Validator.ValidateEmail(user.Email)
	if err != nil {
		return nil, logErrorAndWrap(err,
			"An error occurred when validating the user's email, CreateUser")
	}
	found, err := r.Store.ReadOne(ctx, &ByUsernameOrEmail{user.Username, user.Email})
	if err != nil {
		return nil, logErrorAndWrap(err, "An error occurred on UserStore, CreateUser")
	}
	if found != nil {
		if found.Username != user.Username && found.Email != user.Email {
			return nil, logErrorAndWrap(CorruptedStoreError, fmt.Sprintf(unknownErrorInStore,
				user.Email, user.Username))
		}
		return nil, logErrorAndWrap(EmailOrUsernameAlreadyInUseError, "Existing user")
	}
	err = r.validatePasswords(user.Password, user.RepeatedPassword)
	if err != nil {
		// this error is already wrapped by the validatePasswords function
		return nil, err
	}
	return r.Store.Create(ctx, user)
}

// Updates an user. Returns error on retrieval or actual persistence
func (r *UserRepository) UpdateUser(
	ctx context.Context,
	id string,
	user *UpdateUserDto,
) (*User, error) {
	err := r.Validator.ValidateEmail(user.Email)
	if err != nil {
		return nil, logErrorAndWrap(err,
			"An error occurred when validating the email, UpdateUser")
	}
	found, err := r.Store.ReadOne(ctx, &ByUsernameOrEmail{user.Username, user.Email})
	if err != nil {
		return nil, logErrorAndWrap(UserNotFoundError, "An error occurred on the UserStore, UpdateUser")
	}
	if found != nil {
		if found.Username != user.Username && found.Email != user.Email {
			return nil, logErrorAndWrap(CorruptedStoreError, fmt.Sprintf(unknownErrorInStore,
				user.Email, user.Username))
		}
	}
	r.mapMissingParams(user, found)
	return r.Store.Update(ctx, id, user)
}

func (r *UserRepository) mapMissingParams(user *UpdateUserDto, found *User) {
	if user.Username == "" {
		user.Username = found.Username
	}
	if user.Email == "" {
		user.Email = found.Email
	}
	if user.FirstName == "" {
		user.FirstName = found.FirstName
	}
	if user.LastName == "" {
		user.LastName = found.LastName
	}
}

func (r *UserRepository) validatePasswords(password, repeatedPassword string) error {
	err := r.Validator.ValidatePassword(password)
	if err != nil {
		return logErrorAndWrap(err, "An error occurred on the validator")
	}
	err = r.Validator.ValidatePasswordMatch(password, repeatedPassword)
	if err != nil {
		return logErrorAndWrap(err, "An error occurred on the validator")
	}
	return nil
}

// TODO Remove logErrorAndWrap function as it's unnecessary, users
func logErrorAndWrap(err error, msg string) error {
	return fmt.Errorf("%s: %w \n", msg, err)
}
