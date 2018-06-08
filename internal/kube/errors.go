package kube

import (
	"github.com/pkg/errors"
)

// Kube specific errors
var (
	ErrInvalidID = errors.New("invalid id")
	ErrNotFound  = errors.New("not found")
)
