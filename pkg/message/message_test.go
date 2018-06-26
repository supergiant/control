package message

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestSendUnknownError(t *testing.T) {
	rr := httptest.NewRecorder()
	SendUnknownError(rr, errors.New("test error"))

	msg := new(Message)
	err := jsoniter.NewDecoder(rr.Body).Decode(msg)
	require.NoError(t, err)
	require.Equal(t, sgerrors.UnknownError, msg.ErrorCode)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}
