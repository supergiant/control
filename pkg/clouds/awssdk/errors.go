package awssdk

import "github.com/pkg/errors"

var (
	ErrEmptyRegion  = errors.New("empty region")
	ErrInvalidCreds = errors.New("keyID or secret is not provided")
)
