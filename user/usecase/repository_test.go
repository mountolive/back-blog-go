package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type changePasswordCase struct {
	Name          string
	Description   string
	Repo          *UserRepository
	Dto           *ChangePasswordDto
	ContextCancel bool
	ExpErr        error
}

type createUserCase struct {
	Name          string
	Description   string
	Dto           *CreateUserDto
	Repo          *UserRepository
	ContextCancel bool
	ExpErr        error
}

type updateUserCase struct {
	Name          string
	Description   string
	Dto           *UpdateUserDto
	Repo          *UserRepository
	ContextCancel bool
	ExpErr        error
}

func TestUserRepository(t *testing.T) {
	genericErrMsg := "Got value: %v, Expected: %v"
	happyPathStore := &happyPathUserStoreMock{}
	erroredStore := &erroredUserStoreMock{}
	t.Run("Change Password", func(t *testing.T) {
		correctChangePassword := &ChangePasswordDto{
			Email:            "test@test.com",
			OldPassword:      "Oldpassword1",
			NewPassword:      "Abcdefgh111.",
			RepeatedPassword: "Abcdefgh111.",
		}
		notMatchingPasswords := &ChangePasswordDto{
			Email:            "test@test.com",
			OldPassword:      "Oldpassword1",
			NewPassword:      "Abcdefgh111.",
			RepeatedPassword: "Not matching",
		}
		badEmailChangePassword := &ChangePasswordDto{
			Email:       "test@tes",
			OldPassword: "Oldpassword1",
			NewPassword: "Abcdefgh111.",
		}
		testCases := []changePasswordCase{
			{
				Name:        "Not found user email",
				Description: "It should error out from the UserStore, Email",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     erroredStore,
				},
				Dto:    correctChangePassword,
				ExpErr: UserNotFoundError,
			},
			{
				Name:        "Not found user username",
				Description: "It should error out from the UserStore, Username",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     erroredStore,
				},
				Dto:    correctChangePassword,
				ExpErr: UserNotFoundError,
			},
			{
				Name:        "Malformed email",
				Description: "It should error out from the Validator, Bad Email",
				Repo: &UserRepository{
					Validator: &falseValidatorEmail{},
					Store:     happyPathStore,
				},
				Dto:    badEmailChangePassword,
				ExpErr: MalformedEmailError,
			},
			{
				Name:        "Not valid old password",
				Description: "It should error out from the UserStore, Wrong Old Password",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     erroredStore,
				},
				Dto:    correctChangePassword,
				ExpErr: InvalidOldPasswordError,
			},
			{
				Name:        "Not matching password",
				Description: "It should error out of the validator. Not matching passwords",
				Repo: &UserRepository{
					Validator: &falseValidatorPasswordsNotMatching{},
					Store:     erroredStore,
				},
				Dto:    notMatchingPasswords,
				ExpErr: PasswordsDontMatchError,
			},
			{
				Name:        "Not valid new password",
				Description: "It should error out from the Validator, Bad Password",
				Repo: &UserRepository{
					Validator: &falseValidatorPassword{},
					Store:     happyPathStore,
				},
				Dto:    correctChangePassword,
				ExpErr: InvalidPasswordError,
			},
			{
				Name:        "Context canceled",
				Description: "It should fail with a OperationCanceledError",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				Dto:           correctChangePassword,
				ContextCancel: true,
				ExpErr:        OperationCanceledError,
			},
			{
				Name:        "Valid password",
				Description: "It should change the password properly",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				Dto: correctChangePassword,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.ContextCancel {
					cancel()
				}
				err := tc.Repo.ChangePassword(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErrMsg, err, tc.ExpErr)
				} else {
					require.Nil(t, err, genericErrMsg, err, nil)
				}
			})
		}
	})
	t.Run("Create User", func(t *testing.T) {
		incorrectEmailDto := &CreateUserDto{
			Email:            "non-existence",
			Username:         "correct",
			FirstName:        "Bill",
			LastName:         "Totte",
			Password:         "12345678j.",
			RepeatedPassword: "12345678j.",
		}
		regularDto := &CreateUserDto{
			Email:            "ada@lovelace.com",
			Username:         "ada",
			FirstName:        "Ada",
			LastName:         "Lovelace",
			Password:         "12345678j.",
			RepeatedPassword: "12345678j.",
		}
		notMatchingPasswordDto := &CreateUserDto{
			Email:            "ada@lovelace.com",
			Username:         "ada",
			FirstName:        "Ada",
			LastName:         "Lovelace",
			Password:         "12345678j.",
			RepeatedPassword: "I don't match",
		}
		incorrectPasswordDto := &CreateUserDto{
			Email:     "good@good.com",
			Username:  "good",
			FirstName: "Alexander",
			LastName:  "Grothendiek",
			Password:  "nope",
		}
		happyPathUserStore := &happyPathUserStoreMock{}
		testCases := []createUserCase{
			{
				Name:        "Malformed email",
				Description: "It should fail fast after checking invalid email",
				Repo: &UserRepository{
					Validator: &falseValidatorEmail{},
					Store:     happyPathUserStore,
				},
				Dto:    incorrectEmailDto,
				ExpErr: &InvalidEmailError{errors.New("Bad"), incorrectEmailDto.Email},
			},
			{
				Name:        "Repeated email or username",
				Description: "It should fail after the store finds a matching existing user",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathUserStore,
				},
				Dto:    regularDto,
				ExpErr: UserNotFoundError,
			},
			{
				Name:        "Invalid password",
				Description: "It should fail fast as the password is invalid",
				Repo: &UserRepository{
					Validator: &falseValidatorPassword{},
					Store:     happyPathUserStore,
				},
				Dto:    incorrectPasswordDto,
				ExpErr: InvalidPasswordError,
			},
			{
				Name:        "Context canceled",
				Description: "It should return a proper *UserDto",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				ContextCancel: true,
				Dto:           regularDto,
				ExpErr:        OperationCanceledError,
			},
			{
				Name:        "Not matching passwords",
				Description: "It should return an error from the Validator",
				Repo: &UserRepository{
					Validator: &falseValidatorPasswordsNotMatching{},
					Store:     happyPathStore,
				},
				Dto:    notMatchingPasswordDto,
				ExpErr: PasswordsDontMatchError,
			},
			{
				Name:        "Correct create",
				Description: "It should return a proper *UserDto",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				Dto: regularDto,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.ContextCancel {
					cancel()
				}
				dto, err := tc.Repo.CreateUser(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), tc.Description)
					return
				}
				require.Equal(t, dto.Username, tc.Dto.Username, genericErrMsg,
					dto.Username, tc.Dto.Username)
				require.Equal(t, dto.Email, tc.Dto.Email, genericErrMsg,
					dto.Email, tc.Dto.Email)
				require.Equal(t, dto.FirstName, tc.Dto.FirstName, genericErrMsg,
					dto.FirstName, tc.Dto.FirstName)
				require.Equal(t, dto.LastName, tc.Dto.LastName, genericErrMsg,
					dto.LastName, tc.Dto.LastName)
			})
		}
	})

	t.Run("Update User", func(t *testing.T) {
		incorrectEmailDto := &UpdateUserDto{
			Email:     "bad-company",
			Username:  "correct",
			FirstName: "Samuel",
			LastName:  "Beckett",
		}
		regularDto := &UpdateUserDto{
			Email:     "barbara@liskov.com",
			Username:  "liskov",
			FirstName: "Barb",
			LastName:  "Liskov",
		}
		happyPathUserStore := &happyPathUserStoreMock{}
		testCases := []updateUserCase{
			{
				Name:        "Malformed email",
				Description: "It should fail fast after checking invalid email",
				Repo: &UserRepository{
					Validator: &falseValidatorEmail{},
					Store:     happyPathUserStore,
				},
				Dto:    incorrectEmailDto,
				ExpErr: &InvalidEmailError{errors.New("Bad"), incorrectEmailDto.Email},
			},
			{
				Name:        "Repeated email or username",
				Description: "It should fail after the store finds a matching existing user",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathUserStore,
				},
				Dto:    regularDto,
				ExpErr: UserNotFoundError,
			},
			{
				Name:        "Context canceled",
				Description: "It should fail fast because of the cancelled context",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				ContextCancel: true,
				Dto:           regularDto,
				ExpErr:        OperationCanceledError,
			},
			{
				Name:        "Correct update",
				Description: "It should return a proper *UserDto",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				Dto: regularDto,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.ContextCancel {
					cancel()
				}
				dto, err := tc.Repo.UpdateUser(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), tc.Description)
					return
				}
				require.Equal(t, dto.Username, tc.Dto.Username, genericErrMsg,
					dto.Username, tc.Dto.Username)
				require.Equal(t, dto.Email, tc.Dto.Email, genericErrMsg,
					dto.Email, tc.Dto.Email)
				require.Equal(t, dto.FirstName, tc.Dto.FirstName, genericErrMsg,
					dto.FirstName, tc.Dto.FirstName)
				require.Equal(t, dto.LastName, tc.Dto.LastName, genericErrMsg,
					dto.LastName, tc.Dto.LastName)
			})
		}
	})
}
