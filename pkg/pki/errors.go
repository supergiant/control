package pki

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidCA          = errors.New("certificate is not a certificate authority")
	ErrEmptyPair          = errors.New("pair or cert/key is empty")
	ErrUploadCertificates = errors.New("failed to upload certificates")
)
