package kubernetes

import (
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"k8s.io/client-go/rest"
)

// BuildBasicAuthConfig is a helper function that builds configs for a kubernetes client
// that uses a basic authentication.
// https://kubernetes.io/docs/admin/authentication/#static-password-file
func BuildBasicAuthConfig(host, port, username, pass string) (*rest.Config, error) {
	if host == "" {
		return nil, errors.New("host not specified")
	}
	if username == "" || pass == "" {
		return nil, sgerrors.ErrInvalidCredentials
	}

	if port != "" {
		host += ":" + port
	}

	return &rest.Config{
		Host:     host,
		Username: username,
		Password: pass,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}
