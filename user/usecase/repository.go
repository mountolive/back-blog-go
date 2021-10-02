// This package defines all the regular use cases related to users
package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Id        string
	Email     string
	Username  string
	FirstName string
	LastName  string
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

type LoginDTO struct {
	Email,
	Username,
	Password string
}

// Common errors
var (
	ErrEmailOrUsernameAlreadyInUse = errors.New(
		"email or username passed are already in use")
	ErrUserNotFound = errors.New(
		"user was not found in the DB")
	ErrMalformedEmail = errors.New(
		"email passed is invalid")
	ErrUserPasswordNotMatching = errors.New(
		"seems like user/password data doesn't match with DB record")
	ErrPasswordsDontMatch = errors.New(
		"password and repeatedPassword don't match")
	ErrInvalidPassword = errors.New(
		"password doesn't comply with expected structure")
	ErrOperationCanceled = errors.New(
		"context of the operation was canceled")
	ErrCorruptedStore = errors.New(
		"user' store used is returning inconsistent results")
	ErrCredentialsDontMatch = errors.New("credentials passed don't match")
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

// TODO Split Repository interface by usecase, user service
// Repository defines the basic usecases for the users' domain
type Repository interface {
	ReadUser(ctx context.Context, loginCred string) (*User, error)
	ChangePassword(ctx context.Context, changePass *ChangePasswordDto) error
	CreateUser(ctx context.Context, user *CreateUserDto) (*User, error)
	UpdateUser(ctx context.Context, id string, user *UpdateUserDto) (*User, error)
	Login(ctx context.Context, login LoginDTO) (bool, error)
}

// Basic repository struct. Store is used for persitance and Validator
// for field validation
type UserRepository struct {
	Store     UserStore
	Validator UserValidator
}

var _ Repository = &UserRepository{}

const errMsgLoginRepo = "users repository login: %w"

// Login checks whether the login credentials passed match for the user
func (r *UserRepository) Login(ctx context.Context, login LoginDTO) (bool, error) {
	err := r.Store.CheckIfCorrectPassword(
		ctx,
		&CheckUserAndPasswordDto{
			Email:    login.Email,
			Username: login.Username,
			Password: login.Password,
		},
	)
	if err != nil {
		if errors.Is(err, ErrCredentialsDontMatch) {
			return false, nil
		}
		return false, fmt.Errorf(errMsgLoginRepo, err)
	}
	return true, nil
}

// Reads an user either by her Username or by her Email
func (r *UserRepository) ReadUser(
	ctx context.Context,
	loginCred string,
) (*User, error) {
	if _, err := mail.ParseAddress(loginCred); err != nil {
		return r.Store.ReadOne(ctx, &ByUsernameOrEmail{Username: loginCred})
	}
	return r.Store.ReadOne(ctx, &ByUsernameOrEmail{Email: loginCred})
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
			return nil, logErrorAndWrap(ErrCorruptedStore, fmt.Sprintf(unknownErrorInStore,
				user.Email, user.Username))
		}
		return nil, logErrorAndWrap(ErrEmailOrUsernameAlreadyInUse, "Existing user")
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
		return nil, logErrorAndWrap(ErrUserNotFound, "An error occurred on the UserStore, UpdateUser")
	}
	if found != nil {
		if found.Username != user.Username && found.Email != user.Email {
			return nil, logErrorAndWrap(ErrCorruptedStore, fmt.Sprintf(unknownErrorInStore,
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
	return fmt.Errorf("%s: %w", msg, err)
}
