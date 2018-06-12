package kube

import (
	"github.com/pkg/errors"
)

// Kube specific errors
var (
	ErrNotFound  = errors.New("not found")
)
