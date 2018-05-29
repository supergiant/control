package account

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/asaskevich/govalidator.v8"
	"github.com/gorilla/mux"
	"github.com/supergiant/supergiant/pkg/provider"
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

	req, _ := http.NewRequest(http.MethodPost, "/cloud_accounts", bytes.NewReader(malformedAccount))

	handler := http.HandlerFunc(e.Create)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())

	okAccount, _ := json.Marshal(CloudAccount{
		Name:        "test",
		Provider:    "gce",
		Credentials: nil,
	})
	req, _ = http.NewRequest(http.MethodPost, "/cloud_accounts", bytes.NewReader(okAccount))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	req, _ = http.NewRequest(http.MethodPost, "/cloud_accounts", strings.NewReader("{THIS_IS_INVALID_JSON}"))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())

	m.AssertNumberOfCalls(t, "Put", 1)
}

func TestEndpoint_CreateError(t *testing.T) {
	e, m := fixtures()
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error!"))
	rr := httptest.NewRecorder()

	okAccount, _ := json.Marshal(CloudAccount{
		Name:        "test",
		Provider:    "gce",
		Credentials: nil,
	})
	req, _ := http.NewRequest(http.MethodPost, "/cloud_accounts", bytes.NewReader(okAccount))

	handler := http.HandlerFunc(e.Create)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code, rr.Body.String())
}

func TestEndpoint_Delete(t *testing.T) {
	e, m := fixtures()
	m.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	tt := []struct {
		accountName    string
		responseStatus int
	}{
		{accountName: "NAME", responseStatus: http.StatusOK},
		{accountName: "", responseStatus: http.StatusBadRequest},
	}
	router := mux.NewRouter()
	router.HandleFunc("/cloud_accounts/{accountName}", e.Delete)
	router.HandleFunc("/cloud_accounts/", e.Delete)

	for _, td := range tt {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cloud_accounts/"+td.accountName, nil)

		router.ServeHTTP(rr, req)

		require.Equal(t, td.responseStatus, rr.Code)
	}
}

func TestService_Update(t *testing.T) {
	e, m := fixtures()
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	tt := []struct {
		account        *CloudAccount
		responseStatus int
	}{
		{
			account: &CloudAccount{
				Name:     "OKNAME",
				Provider: provider.AWS,
			},
			responseStatus: http.StatusOK,
		},
		{
			account: &CloudAccount{
				Name:     "NOTOKK",
				Provider: "AAA",
			},
			responseStatus: http.StatusBadRequest,
		},
		{
			account: &CloudAccount{
				Name:     "",
				Provider: provider.DigitalOcean,
			},
			responseStatus: http.StatusBadRequest,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/cloud_accounts/{accountName}", e.Update)
	router.HandleFunc("/cloud_accounts/", e.Update)

	for _, td := range tt {
		body, _ := json.Marshal(td.account)
		m.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(body, nil)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cloud_accounts/"+td.account.Name, bytes.NewReader(body))

		router.ServeHTTP(rr, req)

		require.Equal(t, td.responseStatus, rr.Code)
	}
}
