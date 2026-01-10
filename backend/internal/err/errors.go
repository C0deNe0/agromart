package err

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrUserNotAllowed      = errors.New("user not allowed")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
