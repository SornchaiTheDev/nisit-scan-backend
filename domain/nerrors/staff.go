package nerrors

import "errors"

var (
	ErrStaffAlreadyExists = errors.New("staff already exists")
	ErrStaffNotFound      = errors.New("staff not found")
)
