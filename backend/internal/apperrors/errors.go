package apperrors

import "errors"

var (
	ErrValidation		= errors.New("validation error")
	ErrNotFound			= errors.New("not found")
	ErrConflict			= errors.New("conflict")
	ErrInvalidTimezone	= errors.New("invalid timezone")
)
