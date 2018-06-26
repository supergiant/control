package profile

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodeProfileEndpointCreateProfileError(t *testing.T) {
	endpoint := &NodeProfileHandler{}
	nodeProfile := &NodeProfile{
		ID:           "",
		Size:         &NodeSize{},
		Image:        "ubuntu-18.04",
		Provider:     "openstack",
		Capabilities: []string{},
	}

	data, _ := json.Marshal(nodeProfile)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/nodeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d", http.StatusBadRequest, rr.Code)
	}
}
