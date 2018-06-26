package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestBuildConfigFrom(t *testing.T) {
	tcs := []struct {
		host     string
		port     string
		username string
		pass     string

		config *rest.Config
		err    error
	}{
		// TC#2
		{
			host: "localhost",
			err:  sgerrors.ErrInvalidCredentials,
		},
		// TC#3
		{
			host:     "1.2.3.4",
			username: "user",
			pass:     "pass",
			config: &rest.Config{
				Host:     "1.2.3.4",
				Username: "user",
				Password: "pass",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: true,
				},
			},
		},
		// TC#4
		{
			host:     "1.2.3.4",
			port:     "443",
			username: "user",
			pass:     "pass",
			config: &rest.Config{
				Host:     "1.2.3.4:443",
				Username: "user",
				Password: "pass",
				TLSClientConfig: rest.TLSClientConfig{
					Insecure: true,
				},
			},
		},
	}

	for i, tc := range tcs {
		cfg, err := BuildBasicAuthConfig(tc.host, tc.port, tc.username, tc.pass)
		assert.Equalf(t, tc.err, err, fmt.Sprintf("TC#%d: compare errors", i+1))

		assert.Equalf(t, tc.config, cfg, fmt.Sprintf("TC#%d: compare configs", i+1))
	}
}
