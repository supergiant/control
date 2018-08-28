package sgerrors

type ErrorCode int

const (
	UnknownError        ErrorCode = 1000
	ValidationFailed    ErrorCode = 1001
	InvalidCredentials  ErrorCode = 1003
	NotFound            ErrorCode = 1004
	InvalidJSON         ErrorCode = 1005
	CantChangeID        ErrorCode = 1006
	EntityAlreadyExists ErrorCode = 1007
	TimeoutExceeded     ErrorCode = 1008
)
