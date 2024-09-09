package nerrors

import "errors"

var (
	ErrSomethingWentWrong = errors.New("something went wrong")
	ErrCannotParseUUID    = errors.New("cannot parse uuid")
)
