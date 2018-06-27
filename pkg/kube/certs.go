package kube

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
)

// Certs specific errors.
var (
	ErrInvalidRunner = errors.New("provided invalid runner")

	ErrEmptyName = errors.New("empty name")
)

// Default variables for certs.
const (
	DefaultCertsPath = "/etc/kubernets/ssl"
)

// Bundler represents a key/certificate pair.
type Bundle struct {
	Cert []byte
	Key  []byte
}

// Certs handles keys and certificates for kube components.
type Certs struct {
	path string
	r    runner.Runner
}

// NewCerts returns a configured Certs.
func NewCerts(path string, r runner.Runner) (*Certs, error) {
	if r == nil {
		return nil, ErrInvalidRunner
	}

	return &Certs{
		r:    r,
		path: path,
	}, nil
}

// BundleFor returns keys Bundle for a provided name of the kube component.
func (c *Certs) BundleFor(ctx context.Context, name string) (*Bundle, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	crtBytes, err := c.getFile(ctx, filepath.Join(c.path, certName(name)))
	if err != nil {
		return nil, err
	}

	crtKey, err := c.getFile(ctx, filepath.Join(c.path, keyName(name)))
	if err != nil {
		return nil, err
	}

	return &Bundle{
		Cert: crtBytes,
		Key:  crtKey,
	}, nil
}

func (c *Certs) getFile(ctx context.Context, path string) ([]byte, error) {
	stdout := &bytes.Buffer{}
	cmd := runner.NewCommand(ctx, catCmd(path), stdout)

	err := c.r.Run(cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "run %q", cmd.Script)
	}

	return stdout.Bytes(), nil
}

func catCmd(path string) string {
	return "/usr/bin/cat " + path
}

func keyName(name string) string {
	if strings.TrimSpace(name) == "" {
		return name
	}
	return name + ".key"
}

func certName(name string) string {
	if strings.TrimSpace(name) == "" {
		return name
	}
	return name + ".crt"
}
