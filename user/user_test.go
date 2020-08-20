package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createUserCase struct {
	Name      string
	Username  string
	Email     string
	Password  string
	Validator UserValidator
	Err       error
}

const msg = "Got %v, Expected %v"

// Mocks
type trueValidator struct{}

func (t *trueValidator) ValidateEmail(email string) error {
	return nil
}
func (t *trueValidator) ValidatePassword(password string) error {
	return nil
}

type falseValidatorEmail struct{}

func (f *falseValidatorEmail) ValidateEmail(email string) error {
	return &InvalidEmailError{errors.New("Bad"), email}
}
func (f *falseValidatorEmail) ValidatePassword(password string) error {
	return nil
}

type falseValidatorPassword struct{}

func (f *falseValidatorPassword) ValidateEmail(email string) error {
	return nil
}
func (f *falseValidatorPassword) ValidatePassword(password string) error {
	return InvalidPasswordError
}

// Tests
func TestUser(t *testing.T) {
	tValid := &trueValidator{}
	fValidEmail := &falseValidatorEmail{}
	fValidPass := &falseValidatorPassword{}
	t.Run("NewUser", func(t *testing.T) {
		testCases := []createUserCase{
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
					assert.True(t, invErr.Error() == tcErr.Error(), msg, invErr, tcErr)
				} else if err != nil {
					assert.True(t, errors.Is(err, tc.Err), msg, err, tc.Err)
				} else {
					assert.True(t, user.Username == tc.Username, msg, user.Username, tc.Username)
					assert.True(t, user.Email == tc.Email, msg, user.Email, tc.Email)
					assert.True(t, user.Password == tc.Password, msg, user.Password, tc.Password)
				}
			})
		}
	})
}
