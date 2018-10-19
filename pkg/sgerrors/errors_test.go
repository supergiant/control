package sgerrors

import "testing"

func TestIsNotFound(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			ErrAlreadyExists,
			false,
		},
		{
			ErrNotFound,
			true,
		},
	}

	for _, testCase := range testCases {
		actual := IsNotFound(testCase.err)

		if testCase.expected != actual {
			t.Errorf("Wrong result expected %v actual %v", testCase.expected, actual)
		}
	}
}

func TestIsInvalidCredentials(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			ErrNotFound,
			false,
		},
		{
			ErrInvalidCredentials,
			true,
		},
	}

	for _, testCase := range testCases {
		actual := IsInvalidCredentials(testCase.err)

		if testCase.expected != actual {
			t.Errorf("Wrong result expected %v actual %v", testCase.expected, actual)
		}
	}
}

func TestIsAlreadyExists(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			ErrNotFound,
			false,
		},
		{
			ErrAlreadyExists,
			true,
		},
	}

	for _, testCase := range testCases {
		actual := IsAlreadyExists(testCase.err)

		if testCase.expected != actual {
			t.Errorf("Wrong result expected %v actual %v", testCase.expected, actual)
		}
	}
}

func TestIsUnknownProvider(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			ErrNotFound,
			false,
		},
		{
			ErrUnknownProvider,
			true,
		},
	}

	for _, testCase := range testCases {
		actual := IsUnknownProvider(testCase.err)

		if testCase.expected != actual {
			t.Errorf("Wrong result expected %v actual %v", testCase.expected, actual)
		}
	}
}

func TestIsUnsupportedProvider(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			ErrNotFound,
			false,
		},
		{
			ErrUnsupportedProvider,
			true,
		},
	}

	for _, testCase := range testCases {
		actual := IsUnsupportedProvider(testCase.err)

		if testCase.expected != actual {
			t.Errorf("Wrong result expected %v actual %v", testCase.expected, actual)
		}
	}
}

func TestError_Error(t *testing.T) {
	var (
		code    ErrorCode = 1
		message           = "message"
	)

	err := Error{
		msg:  message,
		Code: code,
	}

	if err.Error() != message {
		t.Errorf("wrong message expected %s actual %s", message, err.Error())
	}
}
