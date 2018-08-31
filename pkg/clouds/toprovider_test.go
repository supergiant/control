package clouds

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToProvider(t *testing.T) {
	tt := []struct {
		str     string
		isValid bool
	}{
		{
			str:     "aws",
			isValid: true,
		},
		{
			str:     "digitalocean",
			isValid: true,
		},
		{
			str:     "openstack",
			isValid: true,
		},
		{
			str:     "packet",
			isValid: true,
		},
		{
			str:     "gce",
			isValid: true,
		},
		{
			str:     "foobar",
			isValid: false,
		},
		{
			str:     "",
			isValid: false,
		},
	}

	for _, tc := range tt {
		_, err := ToProvider(tc.str)
		if tc.isValid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
