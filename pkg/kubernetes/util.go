package kubernetes

import (
	"k8s.io/client-go/rest"
)

// BuildBasicAuthConfig is a helper function that builds configs for a kubernetes client
// that uses a basic authentication.
// https://kubernetes.io/docs/admin/authentication/#static-password-file
func BuildBasicAuthConfig(host, port, username, pass string) (*rest.Config, error) {
	if host == "" {
		return nil, ErrHostNotSpecified
	}
	if username == "" || pass == "" {
		return nil, ErrInvalidCredentials
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
