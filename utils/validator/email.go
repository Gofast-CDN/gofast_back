package validator

import (
	"regexp"
	"unicode/utf8"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type EmailValidator struct {
	MinLength int
	MaxLength int
}

func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		MinLength: 3,
		MaxLength: 254, // RFC 5321
	}
}

func (v *EmailValidator) Validate(email string) error {
	length := utf8.RuneCountInString(email)

	if length < v.MinLength {
		return ErrEmailTooShort
	}

	if length > v.MaxLength {
		return ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}

	return nil
}
