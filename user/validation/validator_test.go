package validation

import (
	"errors"
	"testing"

	"github.com/mountolive/back-blog-go/user/usecase"
	"github.com/stretchr/testify/require"
)

type newValidatorCase struct {
	Name,
	Description,
	EmailRegex,
	PasswordRegex string
	ExpErr error
}

type validateEmailCase struct {
	Name,
	Description,
	Email string
	ExpErr error
}

type validatePasswordCase struct {
	Name,
	Description,
	Password string
	ExpErr error
}

type validateMatchingCase struct {
	Name,
	Description,
	Password,
	Repeated string
	ExpErr error
}

func TestValidator(t *testing.T) {
	genericErr := "Got %s. Expected: %s"
	t.Run("Canary", func(t *testing.T) {
		var _ usecase.UserValidator = &Validator{}
	})

	t.Run("NewValidator", func(t *testing.T) {
		testCases := []newValidatorCase{
			{
				Name:          "Invalid email regex",
				Description:   "Returns an InvalidEmailRegex error and a nil pointer",
				EmailRegex:    "[",
				PasswordRegex: passwordRegex,
				ExpErr:        InvalidEmailRegex,
			},
			{
				Name:          "Invalid password regex",
				Description:   "Returns an InvalidPasswordRegex error and a nil pointer",
				EmailRegex:    emailRegex,
				PasswordRegex: "[",
				ExpErr:        InvalidPasswordRegex,
			},
			{
				Name:          "Valid regex",
				Description:   "Returns a valid Validator pointer",
				EmailRegex:    emailRegex,
				PasswordRegex: passwordRegex,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				validator, err := NewValidator(EmailRegex(tc.EmailRegex), PasswordRegex(tc.PasswordRegex))
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErr, err, tc.ExpErr)
					require.True(t, validator == nil, genericErr, validator, nil)
					return
				}
				require.True(t, err == nil, genericErr, err, tc.ExpErr)
				require.True(t, validator != nil, "Validator shouldn't be nil")
			})
		}
	})

	t.Run("ValidateEmail", func(t *testing.T) {
		validator, err := NewValidator()
		if err != nil {
			t.Errorf("An error occurred at initialization %s \n", err.Error())
		}
		testCases := []validateEmailCase{
			{
				Name:        "Invalid email",
				Description: "Returns an InvalidEmailError for not conforming email",
				Email:       "bad",
				ExpErr:      InvalidEmailError,
			},
			{
				Name:        "Valid email",
				Description: "Returns a nil error, indicating correct email",
				Email:       "abcdef@gh.com",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				err := validator.ValidateEmail(tc.Email)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErr, err, tc.ExpErr)
					return
				}
				require.True(t, err == nil, "Shouldn't return a non nil err: %s", err)
			})
		}
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		validator, err := NewValidator()
		if err != nil {
			t.Errorf("An error occurred at initialization %s \n", err.Error())
		}
		testCases := []validatePasswordCase{
			{
				Name:        "Invalid password short",
				Description: "Returns an InvalidPasswordError for not conforming password: too short",
				Password:    "bad1234",
				ExpErr:      InvalidPasswordError,
			},
			{
				Name:        "Valid password",
				Description: "Returns a nil error, indicating correct password",
				Password:    "AbcdeEfh123093@#$%",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				err := validator.ValidatePassword(tc.Password)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErr, err, tc.ExpErr)
					return
				}
				require.True(t, err == nil, "Shouldn't return a non nil err: %s", err)
			})
		}
	})

	t.Run("ValidatePasswordMatch", func(t *testing.T) {
		validator, err := NewValidator()
		if err != nil {
			t.Errorf("An error occurred at initialization %s \n", err.Error())
		}
		testCases := []validateMatchingCase{
			{
				Name:        "Not matching passwords",
				Description: "Returns an NotMatchingPasswordError",
				Password:    "bad1234",
				Repeated:    "AbcdeEfh123093@#$%",
				ExpErr:      NotMatchingPasswordError,
			},
			{
				Name:        "Matching passwords",
				Description: "Returns a nil error, indicating the passwords do match",
				Password:    "AbcdeEfh123093@#$%",
				Repeated:    "AbcdeEfh123093@#$%",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				err := validator.ValidatePasswordMatch(tc.Password, tc.Repeated)
				if tc.ExpErr != nil {
					require.True(t, errors.Is(err, tc.ExpErr), genericErr, err, tc.ExpErr)
					return
				}
				require.True(t, err == nil, "Shouldn't return a non nil err: %s", err)
			})
		}
	})
}
