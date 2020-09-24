package usecase

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	UserNotFoundError       = errors.New("The user was not found in the DB")
	MalformedEmailError     = errors.New("The email passed is invalid")
	InvalidOldPasswordError = errors.New("The old password passed doesn't match the one expected")
	OperationCanceledError  = errors.New("The context of the operation was canceled")
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
	NewPassword string
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

type Lookup struct {
	FieldName  string
	Comparator Comparator
}

type UserStore interface {
	Create(ctx context.Context, toCreate *User) (*UserDto, error)
	Update(ctx context.Context, updated *User) (*UserDto, error)
	ReadOne(ctx context.Context, criteria ...Lookup) (*UserDto, error)
}

type UserRepository struct {
	Store     UserStore
	Validator UserValidator
}

func (r *UserRepository) ChangePassword(ctx context.Context, changePass *ChangePasswordDto) error {
	// TODO Implement
	return nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *CreateUserDto) (*UserDto, error) {
	// TODO Implement
	return nil, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *UpdateUserDto) (*UserDto, error) {
	// TODO Implement
	return nil, nil
}
