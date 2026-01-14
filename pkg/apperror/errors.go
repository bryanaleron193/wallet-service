package apperror

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotFound      = errors.New("resource not found")
	ErrInsufficient  = errors.New("insufficient balance")
	ErrAlreadyExists = errors.New("resource already exists")
)
