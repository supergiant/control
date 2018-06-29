package node

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNodeCreateError(t *testing.T) {
	endpoint := &Handler{}
	nodeProfile := &Node{
		Id:        "",
		CreatedAt: time.Now().Unix(),
		Provider:  "aws",
		Region:    "us-west1",
	}

	data, _ := json.Marshal(nodeProfile)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/node", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.Put)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d", http.StatusBadRequest, rr.Code)
	}
}
