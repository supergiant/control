package pki

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	val, ok := args.Get(0).(int)
	if !ok {
		return 0, args.Error(1)
	}
	return val, args.Error(1)
}

func TestNewCertificateAuthorityErrPrivateKey(t *testing.T) {
	readErr := errors.New("test")
	m := &mockReader{}
	m.On("Read", mock.Anything).
		Return(mock.Anything, readErr)
	oldReader := rand.Reader
	rand.Reader = m

	_, _, err := newCertificateAuthority()

	if err == nil {
		t.Errorf("error must not be nil")
	}

	if !strings.Contains(err.Error(), readErr.Error()) {
		t.Errorf("Wrong error message %v does not contain %v",
			err, readErr)
	}

	rand.Reader = oldReader
}

func TestNewCertificateAuthoritySuccess(t *testing.T) {
	_, _, err := newCertificateAuthority()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
