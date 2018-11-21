package kube

import (
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
)

var (
	fakeErrFileNotFound = errors.New("key file not found")

	fakeFile = []byte("file")
)

type fakeRunner struct {
	path string
	err  error
}

func newFakeRunner(path string, err error) *fakeRunner {
	return &fakeRunner{path, err}
}

func (r *fakeRunner) Run(cmd *runner.Command) error {
	if cmd.Script == "" {
		return nil
	}

	if !strings.HasSuffix(cmd.Script, r.path) {
		return nil
	}

	if r.err != nil {
		return r.err
	}

	cmd.Out.Write(fakeFile)

	return nil
}

func TestNewCerts(t *testing.T) {
	tcs := []struct {
		path string
		r    runner.Runner

		expected    *Certs
		expectedErr error
	}{
		// TC#1
		{
			expectedErr: ErrInvalidRunner,
		},
		// TC#2
		{
			path:        DefaultCertsPath,
			expectedErr: ErrInvalidRunner,
		},
		// TC#3
		{
			path: DefaultCertsPath,
			r:    &ssh.Runner{},
			expected: &Certs{
				path: DefaultCertsPath,
				r:    &ssh.Runner{},
			},
		},
	}

	for i, tc := range tcs {
		certs, err := NewCerts(tc.path, tc.r)
		require.Equalf(t, tc.expectedErr, err, "TC#%d", i+1)

		require.Equalf(t, tc.expected, certs, "TC#%d", i+1)
	}
}

func TestCerts_BundleFor(t *testing.T) {
	tcs := []struct {
		name string

		runnerFile string
		runnerErr  error

		expected    *Bundle
		expectedErr error
	}{
		// TC#1
		{
			expectedErr: ErrEmptyName,
		},
		// TC#2
		{
			name:        "kubelet",
			runnerFile:  certName("kubelet"),
			runnerErr:   fakeErrFileNotFound,
			expectedErr: fakeErrFileNotFound,
		},
		// TC#3
		{
			name:        "kubelet",
			runnerFile:  keyName("kubelet"),
			runnerErr:   fakeErrFileNotFound,
			expectedErr: fakeErrFileNotFound,
		},
		// TC#4
		{
			name: "kubelet",
			expected: &Bundle{
				Cert: fakeFile,
				Key:  fakeFile,
			},
		},
	}

	for i, tc := range tcs {
		// setup
		certs, err := NewCerts("", newFakeRunner(tc.runnerFile, tc.runnerErr))
		require.Equalf(t, nil, err, "TC#%d: setup certs", i+1)

		// run
		b, err := certs.BundleFor(context.Background(), tc.name)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: get bundle for %s", i+1, tc.name)

		if err == nil {
			require.Equalf(t, tc.expected, b, "TC#%d: compare results", i+1)
		}
	}
}

func TestCerts_getFile(t *testing.T) {
	tcs := []struct {
		path        string
		runnerErr   error
		expectedErr error
	}{
		// TC#1
		{
			path: "kubelet.key",
		},
		// TC#2
		{
			path:        "kubelet.key",
			runnerErr:   fakeErrFileNotFound,
			expectedErr: fakeErrFileNotFound,
		},
	}

	for i, tc := range tcs {
		// setup
		certs, err := NewCerts("", newFakeRunner(tc.path, tc.runnerErr))
		require.Equalf(t, nil, err, "TC#%d: setup certs", i+1)

		// run
		b, err := certs.getFile(context.Background(), tc.path)
		require.EqualValuesf(t, tc.expectedErr, errors.Cause(err), "TC#%d: get file", i+1)

		if err == nil {
			require.Equalf(t, fakeFile, b, "TC#%d: compare results", i+1)
		}
	}
}

func TestKeyName(t *testing.T) {
	tcs := []struct {
		in, out string
	}{
		// TC#1
		{"", ""},
		// TC#1
		{"k", "k.key"},
	}

	for i, tc := range tcs {
		out := keyName(tc.in)
		require.Equalf(t, tc.out, out, "TC#%d", i+i)
	}
}

func TestCertName(t *testing.T) {
	tcs := []struct {
		in, out string
	}{
		// TC#1
		{"", ""},
		// TC#1
		{"k", "k.crt"},
	}

	for i, tc := range tcs {
		out := certName(tc.in)
		require.Equalf(t, tc.out, out, "TC#%d", i+i)
	}
}
