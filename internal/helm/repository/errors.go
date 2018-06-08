package repository

import (
	"github.com/pkg/errors"
)

// Helm repository specific errors
var (
	ErrRepoNil = errors.New("repository is nil")
)
