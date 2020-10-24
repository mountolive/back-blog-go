// This package defines all the regular use cases related to users
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mountolive/back-blog-go/storehelper"
)

type UserDto struct {
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

var basicSearchUserCriteria = []storehelper.Criteria{
	{
		Lookups: []storehelper.Lookup{
			{
				FieldName:  "Email",
				Comparator: storehelper.EQ,
			},
			{
				FieldName:  "Username",
				Comparator: storehelper.EQ,
			},
		},
		Operator: storehelper.OR,
	},
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

var unknownErrorInStore = "Found reported from store: %s and %s, but wrong dto returned"

// Contract for the needs of the repository in terms of persistance:
//    Defines which methods would be needed for each usecase
type UserStore interface {
	Create(context.Context, *CreateUserDto) (*UserDto, error)
	Update(context.Context, *UpdateUserDto) (*UserDto, error)
	UpdatePassword(context.Context, *ChangePasswordDto) error
	ReadOne(context.Context, ...storehelper.Criteria) (*UserDto, error)
	CheckIfCorrectPassword(context.Context, *CheckUserAndPasswordDto) error
}

type Logger interface {
	LogError(err error)
}

// Contract for the needs of the repository in terms of validation:
//     Methods needed by each usecase for validating the user's data
type UserValidator interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidatePasswordMatch(password, repeatedPassword string) error
}

// Basic repository struct. Store is used for persitance and Validator
// for field validation
type UserRepository struct {
	Store     UserStore
	Validator UserValidator
	Logger    Logger
}

// Changes password and persists. Returns an error on validation or
// store's retrieval/persistence
func (r *UserRepository) ChangePassword(ctx context.Context,
	changePass *ChangePasswordDto) error {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		// Already wrapped by the helper function
		return err
	}
	defer cancel()
	err = r.Validator.ValidateEmail(changePass.Email)
	if err != nil {
		return r.logErrorAndWrap(err,
			"An error occurred when validating the user's email, ChangePassword")
	}
	toCheck := &CheckUserAndPasswordDto{
		Email:    changePass.Email,
		Username: changePass.Username,
		Password: changePass.OldPassword,
	}
	err = r.Store.CheckIfCorrectPassword(ctx, toCheck)
	if err != nil {
		return r.logErrorAndWrap(err, "An error occurred on the UserStore, ChangePassword")
	}
	err = r.validatePasswords(changePass.NewPassword, changePass.RepeatedPassword)
	if err != nil {
		// This error is already wrapped by the validatePasswords function
		return err
	}
	return r.Store.UpdatePassword(ctx, changePass)
}

// Creates an user. Returns an error on validation
func (r *UserRepository) CreateUser(ctx context.Context,
	user *CreateUserDto) (*UserDto, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		// Already wrapped by the helper function
		return nil, err
	}
	defer cancel()
	err = r.Validator.ValidateEmail(user.Email)
	if err != nil {
		return nil, r.logErrorAndWrap(err,
			"An error occurred when validating the user's email, CreateUser")
	}
	found, err := r.Store.ReadOne(ctx, basicSearchUserCriteria...)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "An error occurred on the UserStore, CreateUser")
	}
	if found != nil {
		if found.Username != user.Username && found.Email != user.Email {
			return nil, r.logErrorAndWrap(CorruptedStoreError, fmt.Sprintf(unknownErrorInStore,
				user.Email, user.Username))
		}
		return nil, r.logErrorAndWrap(EmailOrUsernameAlreadyInUseError, "Existing user")
	}
	err = r.validatePasswords(user.Password, user.RepeatedPassword)
	if err != nil {
		// this error is already wrapped by the validatePasswords function
		return nil, err
	}
	return r.Store.Create(ctx, user)
}

// Updates an user. Returns error on retrieval or actual persistence
func (r *UserRepository) UpdateUser(ctx context.Context,
	user *UpdateUserDto) (*UserDto, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		// Already wrapped by the helper function
		return nil, err
	}
	defer cancel()
	err = r.Validator.ValidateEmail(user.Email)
	if err != nil {
		return nil, r.logErrorAndWrap(err,
			"An error occurred when validating the email, UpdateUser")
	}
	found, err := r.Store.ReadOne(ctx, basicSearchUserCriteria...)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "An error occurred on the UserStore, UpdateUser")
	}
	if found != nil {
		if found.Username != user.Username && found.Email != user.Email {
			return nil, r.logErrorAndWrap(CorruptedStoreError, fmt.Sprintf(unknownErrorInStore,
				user.Email, user.Username))
		}
	}
	r.mapMissingParams(user, found)
	return r.Store.Update(ctx, user)
}

func (r *UserRepository) mapMissingParams(user *UpdateUserDto, found *UserDto) {
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
		return r.logErrorAndWrap(err, "An error occurred on the validator")
	}
	err = r.Validator.ValidatePasswordMatch(password, repeatedPassword)
	if err != nil {
		return r.logErrorAndWrap(err, "An error occurred on the validator")
	}
	return nil
}

func (r *UserRepository) logErrorAndWrap(err error, msg string) error {
	r.Logger.LogError(err)
	return fmt.Errorf("%s: %w \n", msg, err)
}

func checkContextAndRecreate(ctx context.Context) (context.Context, context.CancelFunc, error) {
	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("%w: %v", OperationCanceledError, ctx.Err())
	default:
		newCtx, cancel := context.WithCancel(ctx)
		return newCtx, cancel, nil
	}
}
