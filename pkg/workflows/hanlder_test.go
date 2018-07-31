package workflows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/gorilla/mux"

	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

func TestWorkflowHandlerGetWorkflow(t *testing.T) {
	id := "abcd"
	expectedType := "master"
	expectedSteps := []StepStatus{{}, {}}
	w1 := &WorkFlow{
		Type:         expectedType,
		StepStatuses: expectedSteps,
	}
	data, _ := json.Marshal(w1)

	h := WorkflowHandler{
		repository: &fakeRepository{
			map[string][]byte{
				fmt.Sprintf("workflows/%s", id): data,
			},
		},
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/workflows/%s", id), nil)

	router := mux.NewRouter()
	router.HandleFunc("/workflows/{id}", h.GetWorkflow)
	router.ServeHTTP(resp, req)

	w2 := &WorkFlow{}
	err := json.Unmarshal(resp.Body.Bytes(), w2)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if len(w1.StepStatuses) != len(w2.StepStatuses) {
		t.Errorf("Wrong step statuses len expected  %d actual %d",
			len(w1.StepStatuses), len(w2.StepStatuses))
	}

	if w1.Type != w2.Type {
		t.Errorf("Wrong workflow type expected %s actual %s",
			w1.Type, w2.Type)
	}
}

func TestWorkflowHandlerBuildWorkflow(t *testing.T) {
	h := WorkflowHandler{
		repository: &fakeRepository{
			map[string][]byte{},
		},
	}

	reqBody := BuildWorkFlowRequest{
		SshConfig: ssh.Config{
			Host:    "12.34.56.67",
			Port:    "22",
			User:    "root",
			Timeout: 1,
			Key:     []byte("")},
	}
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(reqBody)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", body)

	h.BuildWorkflow(rec, req)

	resp := &BuildWorkFlowResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), resp)

	if rec.Code != http.StatusCreated {
		t.Errorf("Wrong response code expected %d actual %d", rec.Code, http.StatusCreated)
		return
	}

	if err != nil {
		t.Errorf("Unexpected error while parsing response %v", err)
	}
}
