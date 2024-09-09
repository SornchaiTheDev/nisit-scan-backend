package nerrors

import "errors"

var (
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminAlreadyExists = errors.New("admin already exists")
)
