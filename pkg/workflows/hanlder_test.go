package workflows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
