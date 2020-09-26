package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Common errors
var (
	EmailOrUsernameAlreadyInUseError = errors.New("The email or username passed are already in use")
	UserNotFoundError                = errors.New("The user was not found in the DB")
	MalformedEmailError              = errors.New("The email passed is invalid")
	InvalidOldPasswordError          = errors.New("The old password passed doesn't match the one expected")
	PasswordsDontMatchError          = errors.New("Password and RepeatedPassword don't match")
	InvalidPasswordError             = errors.New("Password doesn't comply with expected structure")

	OperationCanceledError = errors.New("The context of the operation was canceled")
	CorruptedStoreError    = errors.New("The UserStore used is returning inconsistent results")
)

type Post struct {
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserDto struct {
	Id        string
	Email     string
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post
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

// Single lookup condition
type Lookup struct {
	FieldName  string
	Comparator Comparator
}

// Composed search criteria
// Unifies Lookups either by AND or by ORs
type Criteria struct {
	Lookups  []Lookup
	Operator LogicalOperator
}

type InvalidEmailError struct {
	Err      error
	BadEmail string
}

type Comparator int

// EQ   ->  =
// NEQ  ->  !=
// LET  ->  <=
// LT   ->  <
// GET  ->  >=
// GT   ->  >
const (
	EQ Comparator = iota
	NEQ
	LET
	LT
	GET
	GT
)

type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
)

var basicSearchUserCriteria = []Criteria{
	{
		Lookups: []Lookup{
			{
				FieldName:  "Email",
				Comparator: EQ,
			},
			{
				FieldName:  "Username",
				Comparator: EQ,
			},
		},
		Operator: OR,
	},
}

var unknownErrorInStore = "Found reported from store: %s and %s, but wrong dto returned"

// Contract for the needs of the repository in terms of persistance:
//    Defines which methods would be needed for each usecase
type UserStore interface {
	Create(context.Context, *CreateUserDto) (*UserDto, error)
	Update(context.Context, *UpdateUserDto) (*UserDto, error)
	UpdatePassword(context.Context, *ChangePasswordDto) error
	ReadOne(context.Context, ...Criteria) (*UserDto, error)
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

func (u *InvalidEmailError) Error() string {
	return fmt.Sprintf("The passed email: %v, has an error: %v \n",
		u.BadEmail, u.Err)
}

func (u *InvalidEmailError) Unwrap() error {
	return u.Err
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
