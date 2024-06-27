package interaction

import "errors"

var (
	ErrNotAllowed                = errors.New("user is not allowed to perform this action")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrUserNotSubscribed         = errors.New("user is not subscribed")
)
