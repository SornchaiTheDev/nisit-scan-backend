package nerrors

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventAlreadyExists = errors.New("event already exists")
)
