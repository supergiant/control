package ssh

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/runner"
)

func TestRunner_New(t *testing.T) {
	testCases := []struct {
		conf           Config
		cmd            *runner.Command
		expectedRunner *Runner
		expectedErr    error
	}{
		{
			expectedErr: ErrHostNotSpecified,
		},
		{
			conf: Config{
				Host: "localhost",
			},
			expectedErr: ErrUserNotSpecified,
		},
		{
			conf: Config{
				Host: "localhost",
				User: "sg",
			},
			expectedErr: errors.New("ssh: no key found"),
		},
		{
			conf: Config{
				Host: "localhost",
				User: "sg",
				Key:  []byte("12345"),
			},
			expectedErr: errors.New("ssh: no key found"),
		},
		{
			conf: Config{
				Host: "localhost",
				User: "sg",
				Key:  []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC0/5OuS6CYjtNR1y8Nzn4Uzn+s7XAl/E4T7LfQ5Rf/Wf+23zWP\nIi+NCLF9sD/ij2UPZfONXrsdAk0bQsqyWycHE3dYxXoR2yROtq5VfnG1IJz4y/wx\nt4iuS5Lk6wCEiIsu9bH1jNVKP4vZ9QRTutEV2KTDTz+pNFuD6O80fgpgUQIDAQAB\nAoGAZH7TXJ8ZGBuVMYes1JmmX58YPMfI0Q55u45fsVjCLkGmPb9JWaO9iy0cf5Dw\np7a+ggI1bHUAp2azsUMwkD8UN2aJsf8rfYXo3OgUqbuYRo6n8Pf3V9Hb2ciNTE9Q\nnpE0eoNf6LoN0aIfl+qH7GxDP2Ql/+WKzbwhpBYlP+tj7VECQQDX8hei8hiyThcd\n3yB0I1UmcGW5TMjUtBF2w8HCpYEj8nqi6EFrysSVjA/i3vtq5pqgO9+yLbgujhif\nl12kj6BLAkEA1pIODXiadGg4xa75CqE1EcuL5cxBghVlY/EDYcj5wX9ES2XbroRR\nt0XuTT5QQzLOd1LTfe1Oiw4dMkILIZs4UwJBAMPSr4h+DdMzaVcXTXDD0aWn6zcb\n4Eqyd9vBLOX7+53Dd15fS2QaXiZW+tj65/dK4xFG+lWzi//7r2yZcLuX2v0CQB+t\nIiv08QBcXn04joV2NQpyfS2okMcud3BgpTorXEuniSKEYAEMga/HwB1hJKI2/un4\nrUY64UyAAelofJIygwcCQFNySN/ZRoDgIEDq69k6bz4Z+CTJAdJ07OudP3EzRFor\nYHZAiFRTUUdcsX2iNmtZy/b0IqxfSGW/T9PmwaF6ZV0=\n-----END RSA PRIVATE KEY-----\n"),
			},
		},
	}

	for i, tc := range testCases {
		_, err := NewRunner(tc.conf)
		if tc.expectedErr != errors.Cause(err) {
			// ssh errors are hardcoded, so check the msg too..
			if tc.expectedErr.Error() != errors.Cause(err).Error() {
				require.Equalf(t, tc.expectedErr, err, "TC#%d", i+1)
			}
		}
	}
}
