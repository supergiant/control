package kube

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKubeProfileEndpointCreateProfile(t *testing.T) {
	endpoint := &KubeEndpoint{}
	kube := &Kube{
		Name:              "",
		KubernetesVersion: "1.8.7",
		MasterPublicIP:    "12.34.56.78",
		RBACEnabled:       false,
	}

	data, _ := json.Marshal(kube)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/kube", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateKube)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d", http.StatusBadRequest, rr.Code)
	}
}
