package nerrors

import "errors"

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenNotValid = errors.New("token not valid")
	ErrTokenNotMatch = errors.New("token not match")
)
