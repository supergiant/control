package account

import (
	"testing"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/asaskevich/govalidator.v8"
	"strings"
	"net/http/httptest"
	"github.com/pkg/errors"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	args := m.Called(ctx, prefix)
	return args.Get(0).([][]byte), args.Error(1)
}

func (m *MockStorage) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	args := m.Called(ctx, prefix, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorage) Put(ctx context.Context, prefix string, key string, value []byte) error {
	args := m.Called(ctx, prefix, key, value)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, prefix string, key string) error {
	args := m.Called(ctx, prefix, key)
	return args.Error(0)
}

func fixtures() (*Endpoint, *MockStorage) {
	mockStorage := new(MockStorage)
	return &Endpoint{
		Service: &Service{
			Repository: mockStorage,
		},
	}, mockStorage
}
func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func TestEndpoint_Create(t *testing.T) {
	e, m := fixtures()
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	malformedAccount, _ := json.Marshal(CloudAccount{
		Name:        "",
		Provider:    "asdasd",
		Credentials: nil,
	})

	req, _ := http.NewRequest("POST", "/cloud_accounts", bytes.NewReader(malformedAccount))

	handler := http.HandlerFunc(e.Create)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())

	okAccount, _ := json.Marshal(CloudAccount{
		Name:        "test",
		Provider:    "gce",
		Credentials: nil,
	})
	req, _ = http.NewRequest("POST", "/cloud_accounts", bytes.NewReader(okAccount))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	req, _ = http.NewRequest("POST", "/cloud_accounts", strings.NewReader("{THIS_IS_INVALID_JSON}"))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())

	m.AssertNumberOfCalls(t, "Put", 1)
}

func TestEndpoint_CreateError(t *testing.T) {
	e, m := fixtures()
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("ERROR!"))
	rr := httptest.NewRecorder()

	okAccount, _ := json.Marshal(CloudAccount{
		Name:        "test",
		Provider:    "gce",
		Credentials: nil,
	})
	req, _ := http.NewRequest("POST", "/cloud_accounts", bytes.NewReader(okAccount))

	handler := http.HandlerFunc(e.Create)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code, rr.Body.String())
}
