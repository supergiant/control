package amazon

import "github.com/pkg/errors"

var (
	ErrReadVPC       = errors.New("aws: can't read vpc info")
	ErrCreateVPC     = errors.New("aws: create vpc")
	ErrAuthorization = errors.New("aws: authorization")
)
