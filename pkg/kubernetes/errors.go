package kubernetes

import (
	"errors"
)

var (
	ErrHostNotSpecified = errors.New("host not specified")

	ErrInvalidCredentials = errors.New("invalid credentials: user or password not specified")
)
