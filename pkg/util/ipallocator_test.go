package util

import (
	"net"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestGetDNSIP(t *testing.T) {
	for _, tc := range []struct {
		name        string
		subnet      string
		expectedIP  string
		expectedErr error
	}{
		{
			name:       "success",
			subnet:     "10.3.0.0/16",
			expectedIP: "10.3.0.10",
		},
		{
			name:        "small_subnet",
			subnet:      "10.3.0.0/32",
			expectedErr: ErrSubnetTooSmall,
		},
		{
			name:        "invalid_subnet",
			subnet:      "10.3.0.0/126",
			expectedErr: &net.ParseError{Type: "CIDR address", Text: "10.3.0.0/126"},
		},
	} {
		ip, err := GetDNSIP(tc.subnet)
		if !reflect.DeepEqual(errors.Cause(err), tc.expectedErr) {
			t.Fatalf("TC: %s: error: %q, expected error: %q", tc.name, errors.Cause(err), tc.expectedErr)

		}
		if err == nil && ip.String() != tc.expectedIP {
			t.Fatalf("TC: %s: subnet %s, ip %s, expected ip %s", tc.name, tc.subnet, ip, tc.expectedIP)
		}
	}
}

func TestGetKubernetesDefaultSvcIP(t *testing.T) {
	for _, tc := range []struct {
		name       string
		subnet     string
		expectedIP string
	}{
		{
			name:       "success",
			subnet:     "10.3.0.0/16",
			expectedIP: "10.3.0.1",
		},
	} {
		ip, err := GetKubernetesDefaultSvcIP(tc.subnet)
		if err != nil {
			t.Fatalf("TC: %s: error: %q, expected error: <nil>", tc.name, errors.Cause(err))
		}
		if err == nil && ip.String() != tc.expectedIP {
			t.Fatalf("TC: %s: subnet %s, ip %s, expected ip %s", tc.name, tc.subnet, ip, tc.expectedIP)
		}
	}
}
