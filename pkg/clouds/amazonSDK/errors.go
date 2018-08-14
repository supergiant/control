package amazonSDK

import "github.com/pkg/errors"

var (
	ErrInvalidCreds = errors.New("keyID or secret is not provided")
)
