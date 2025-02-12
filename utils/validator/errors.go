package validator

import "errors"

var (
	ErrEmailTooShort      = errors.New("email is too short")
	ErrEmailTooLong       = errors.New("email exceeds maximum length")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)
