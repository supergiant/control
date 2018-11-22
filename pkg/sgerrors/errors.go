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
	ErrInvalidCredentials  = New("invalid credentials", InvalidCredentials)
	ErrNotFound            = New("entity not found", NotFound)
	ErrAlreadyExists       = New("entity already exists", EntityAlreadyExists)
	ErrUnknownProvider     = New("unknown provider type", UnknownProvider)
	ErrUnsupportedProvider = New("unsupported provider", UnsupportedProvider)
	ErrInvalidJson         = New("invalid json", InvalidJSON)
	ErrNilValue            = New("nil value", NilValue)
	ErrTokenExpired        = New("token has been expire", TokenExpired)
	ErrNilEntity           = New("nil entity", NilEntity)
	ErrTimeoutExceeded     = New("timeout exceeded", TimeoutExceeded)
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

func IsTimeoutExceeded(err error) bool {
	return errors.Cause(err) == ErrTimeoutExceeded
}

func IsUnknownProvider(err error) bool {
	return errors.Cause(err) == ErrUnknownProvider
}

func IsUnsupportedProvider(err error) bool {
	return errors.Cause(err) == ErrUnsupportedProvider
}
