package profile

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/asaskevich/govalidator.v8"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func TestKubeProfileEndpointCreateProfile(t *testing.T) {
	endpoint := &KubeProfileHandler{}
	kubeProfile := &Profile{
		ID:          "",
		K8SVersion:  "1.11.1",
		RBACEnabled: false,
	}

	data, _ := json.Marshal(kubeProfile)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/kubeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d", http.StatusBadRequest, rr.Code)
	}
}
