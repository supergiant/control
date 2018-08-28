package pki

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidCA          = errors.New("certificate is not a certificate authority")
	ErrUploadCertificates = errors.New("failed to upload certificates")
)
