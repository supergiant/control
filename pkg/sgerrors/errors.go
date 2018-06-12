package sgerrors

import (
	"github.com/pkg/errors"
)

type Error struct {
	msg  string
	Code ErrorCode
}

func (e *Error) Error() string {
	return e.msg
}

func New(msg string, code ErrorCode) error {
	return &Error{
		msg:  msg,
		Code: code,
	}
}

var (
	ErrInvalidCredentials = New("invalid credentials", InvalidCredentials)
	ErrNotFound           = New("entity not found", NotFound)
	ErrAlreadyExists      = New("entity already exists", EntityAlreadyExists)
)

func IsNotFound(err error) bool {
	return errors.Cause(err) == ErrNotFound
}

func IsInvalidCredentials(err error) bool {
	return errors.Cause(err) == ErrInvalidCredentials
}

func IsAlreadyExists(err error) bool {
	return errors.Cause(err) == ErrAlreadyExists
}
