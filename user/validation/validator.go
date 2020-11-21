// Holds the logic associated to validation of users' fields
package validation

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	emailRegex    = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	passwordRegex = `[A-Za-z0-9@!.,~#$%&*<>+"]{8,}$`
)

// Used to validate whether an email or a password
// string are correct
// Used also to validate equality between 2 passwords' strings
type Validator struct {
	emailRegex,
	passwordRegex *regexp.Regexp
}

var (
	InvalidEmailRegex    = errors.New("The configured email regex couldn't be compiled")
	InvalidPasswordRegex = errors.New("The configured password regex couldn't be compiled")

	InvalidPasswordError     = errors.New("The password passed doesn't comply the min requirements")
	InvalidEmailError        = errors.New("The email passed is not a valid email")
	NotMatchingPasswordError = errors.New("The passwords are not matching")
)

// Creates a new email and password validator
// that will use the passed email and password regexes
func NewValidator(regexes ...func(*Validator) error) (*Validator, error) {
	validator := &Validator{}
	for _, regexer := range regexes {
		err := regexer(validator)
		if err != nil {
			return nil, err
		}
	}
	if validator.emailRegex == nil {
		err := validator.setEmailRegex(emailRegex)
		if err != nil {
			return nil, err
		}
	}
	if validator.passwordRegex == nil {
		err := validator.setPasswordRegex(passwordRegex)
		if err != nil {
			return nil, err
		}
	}
	return validator, nil
}

// Validates an email address using the internal regexp
// representation of an email, holded by the Validator instance
func (v *Validator) ValidateEmail(email string) error {
	matches := v.emailRegex.MatchString(email)
	if !matches {
		return InvalidEmailError
	}
	return nil
}

// Validates a password using the internal regexp
// representation of a password, holded by the Validator instance
func (v *Validator) ValidatePassword(password string) error {
	matches := v.passwordRegex.MatchString(password)
	if !matches {
		return InvalidPasswordError
	}
	return nil
}

// Checks whether 2 password's strings are the same
func (v *Validator) ValidatePasswordMatch(pass, repeat string) error {
	if pass != repeat {
		return NotMatchingPasswordError
	}
	return nil
}

// Assigns a custom emailRegex to the validator instance
func EmailRegex(regex string) func(*Validator) error {
	return func(v *Validator) error {
		return v.setEmailRegex(regex)
	}
}

// Assigns a custom passwordRegex to the validator instance
func PasswordRegex(regex string) func(*Validator) error {
	return func(v *Validator) error {
		return v.setPasswordRegex(regex)
	}
}

func (v *Validator) setEmailRegex(regex string) error {
	compEmail, err := regexp.Compile(regex)
	if err != nil {
		return wrapError(InvalidEmailRegex, err)
	}
	v.emailRegex = compEmail
	return nil
}

func (v *Validator) setPasswordRegex(regex string) error {
	compPassword, err := regexp.Compile(regex)
	if err != nil {
		return wrapError(InvalidPasswordRegex, err)
	}
	v.passwordRegex = compPassword
	return nil
}

func wrapError(wrapperError, err error) error {
	return fmt.Errorf("%w: %s", wrapperError, err.Error())
}
