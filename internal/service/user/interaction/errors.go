package interaction

import "errors"

var (
	ErrNotAllowed = errors.New("user is not allowed to perform this action")
)
