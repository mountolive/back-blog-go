package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type newUserCase struct {
	Name      string
	Username  string
	Email     string
	Password  string
	Validator UserValidator
	Err       error
}

func TestUser(t *testing.T) {
	genericErrMsg := "Got result: %v, but Expected: %v"
	tValid := &trueValidator{}
	fValidEmail := &falseValidatorEmail{}
	fValidPass := &falseValidatorPassword{}
	t.Run("NewUser", func(t *testing.T) {
		testCases := []newUserCase{
			{
				Name:      "Correct email and password",
				Username:  "any",
				Email:     "test@gmail.com",
				Password:  "abcdefg1",
				Validator: tValid,
				Err:       nil,
			},
			{
				Name:      "Invalid email",
				Username:  "any",
				Email:     "invalid@invalid",
				Password:  "abcdefg1",
				Validator: fValidEmail,
				Err:       &InvalidEmailError{errors.New("Bad"), "invalid@invalid"},
			},
			{
				Name:      "Invalid password",
				Username:  "any",
				Email:     "test@gmail.com",
				Password:  "1a",
				Validator: fValidPass,
				Err:       InvalidPasswordError,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				user, err := NewUser(tc.Username, tc.Email, tc.Password, tc.Validator)
				var invalid *InvalidEmailError
				if err != nil && errors.As(err, &invalid) {
					invErr := invalid.Unwrap()
					tcErr := errors.Unwrap(tc.Err)
					assert.True(t, invErr.Error() == tcErr.Error(), genericErrMsg,
						invErr, tcErr)
				} else if err != nil {
					assert.True(t, errors.Is(err, tc.Err), genericErrMsg, err, tc.Err)
				} else {
					assert.True(t, user.Username == tc.Username, genericErrMsg,
						user.Username, tc.Username)
					assert.True(t, user.Email == tc.Email, genericErrMsg,
						user.Email, tc.Email)
					assert.True(t, user.Password == tc.Password, genericErrMsg,
						user.Password, tc.Password)
				}
			})
		}
	})
}
