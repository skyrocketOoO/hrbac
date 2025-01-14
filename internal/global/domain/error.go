package dm

import "errors"

var (
	ErrNotImplement     = errors.New("not implemented")
	ErrEmptyRequest     = errors.New("empty request body")
	ErrUnknown          = errors.New("unknown")
	ErrLoginFailed      = errors.New("login failed")
	ErrUserNameRepetite = errors.New("user name repetite")
)
