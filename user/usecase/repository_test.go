package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type changePasswordCase struct {
	Name   string
	Repo   *UserRepository
	Dto    *ChangePasswordDto
	ExpErr error
}

type createUserCase struct {
	Name   string
	Dto    *CreateUserDto
	ExpErr error
}

type updateUserCase struct {
	Name   string
	Dto    *UpdateUserDto
	ExpErr error
}

func TestUserRepository(t *testing.T) {
	genericErrMsg := "Got value: %v, Expected: %v"
	correctChangePassword := &ChangePasswordDto{
		Email:       "test@test.com",
		OldPassword: "oldpassword1",
		NewPassword: "Abcdefgh111.",
	}
	badEmailChangePassword := &ChangePasswordDto{
		Email:       "test@tes",
		OldPassword: "oldpassword1",
		NewPassword: "Abcdefgh111.",
	}
	happyPathStore := &happyPathUserStoreMock{}
	erroredStore := &erroredUserStoreMock{}
	t.Run("Change Password", func(t *testing.T) {
		testCases := []changePasswordCase{
			{
				Name: "Not found user",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     erroredStore,
				},
				Dto:    correctChangePassword,
				ExpErr: UserNotFoundError,
			},
			{
				Name: "Malformed email",
				Repo: &UserRepository{
					Validator: &falseValidatorEmail{},
					Store:     happyPathStore,
				},
				Dto:    badEmailChangePassword,
				ExpErr: MalformedEmailError,
			},
			{
				Name: "Not valid old password",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     erroredStore,
				},
				Dto:    correctChangePassword,
				ExpErr: InvalidOldPasswordError,
			},
			{
				Name: "Not valid new password",
				Repo: &UserRepository{
					Validator: &falseValidatorPassword{},
					Store:     happyPathStore,
				},
				Dto:    correctChangePassword,
				ExpErr: InvalidPasswordError,
			},
			{
				Name: "Valid password",
				Repo: &UserRepository{
					Validator: &trueValidator{},
					Store:     happyPathStore,
				},
				Dto: correctChangePassword,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				err := tc.Repo.ChangePassword(context.Background(), tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErrMsg, err, tc.ExpErr)
				} else {
					require.Nil(t, err, genericErrMsg, err, nil)
				}
			})
		}
	})
}
