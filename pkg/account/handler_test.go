package account

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/testutils"
)

func fixtures() (*Handler, *testutils.MockStorage) {
	mockStorage := new(testutils.MockStorage)
	return &Handler{
		service: &Service{
			storagePrefix: DefaultStoragePrefix,
			repository:    mockStorage,
		},
	}, mockStorage
}
func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func TestEndpoint_Create(t *testing.T) {
	e, m := fixtures()
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	malformedAccount, _ := json.Marshal(model.CloudAccount{
		Name:        "ff",
		Provider:    clouds.DigitalOcean,
		Credentials: map[string]string{},
	})

	req, _ := http.NewRequest(http.MethodPost, "/cloud_accounts", bytes.NewReader(malformedAccount))

	handler := http.HandlerFunc(e.Create)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code, rr.Body.String())

	okAccount, _ := json.Marshal(model.CloudAccount{
		Name:     "test",
		Provider: clouds.DigitalOcean,
		Credentials: map[string]string{
			clouds.DigitalOceanAccessToken: "test",
			clouds.DigitalOceanFingerPrint: "fingerprint",
		},
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

	okAccount, _ := json.Marshal(model.CloudAccount{
		Name:        "test",
		Provider:    "gce",
		Credentials: map[string]string{},
	})
	req, _ := http.NewRequest(http.MethodPost, "/cloud_accounts", bytes.NewReader(okAccount))

	handler := http.HandlerFunc(e.Create)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())
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
		account        *model.CloudAccount
		responseStatus int
	}{
		{
			account: &model.CloudAccount{
				Name:     "OKNAME",
				Provider: clouds.AWS,
			},
			responseStatus: http.StatusOK,
		},
		{
			account: &model.CloudAccount{
				Name:     "NOTOKK",
				Provider: "AAA",
			},
			responseStatus: http.StatusBadRequest,
		},
		{
			account: &model.CloudAccount{
				Name:     "",
				Provider: clouds.DigitalOcean,
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
