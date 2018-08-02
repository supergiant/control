package workflows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bytes"

	"github.com/gorilla/mux"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestWorkflowHandlerGetWorkflow(t *testing.T) {
	id := "abcd"
	expectedType := "master"
	expectedSteps := []StepStatus{{}, {}}
	w1 := &Task{
		Type:         expectedType,
		StepStatuses: expectedSteps,
	}
	data, _ := json.Marshal(w1)

	h := TaskHandler{
		repository: &mockRepository{
			map[string][]byte{
				fmt.Sprintf("%s/%s", prefix, id): data,
			},
		},
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", prefix, id), nil)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s/{id}", prefix), h.GetTask)
	router.ServeHTTP(resp, req)

	w2 := &Task{}
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

func TestTaskHandlerRunTask(t *testing.T) {
	h := TaskHandler{
		runnerFactory: func(cfg ssh.Config) (runner.Runner, error) {
			return &testutils.FakeRunner{}, nil
		},
		repository: &mockRepository{
			map[string][]byte{},
		},
	}

	workflowName := "workflow1"
	message := "hello, world!!!"
	step := &mockStep{
		name:     "mock_step",
		messages: []string{message},
	}

	workflow := []steps.Step{step}
	RegisterWorkFlow(workflowName, workflow)

	reqBody := RunTaskRequest{
		Cfg: steps.Config{
			Timeout: time.Second * 1,
		},
		WorkflowName: workflowName,
	}

	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(reqBody)

	if err != nil {
		t.Errorf("Unexpected error while json encoding req body %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	h.RunTask(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Wrong response code expected %d a ctual %d", http.StatusAccepted, rec.Code)
	}

	resp := &TaskResponse{}
	json.NewDecoder(rec.Body).Decode(resp)

	if len(resp.Id) == 0 {
		t.Error("task id in response should not be empty")
	}
}

func TestWorkflowHandlerBuildWorkflow(t *testing.T) {
	h := TaskHandler{
		runnerFactory: func(cfg ssh.Config) (runner.Runner, error) {
			return &testutils.FakeRunner{}, nil
		},
		repository: &mockRepository{
			map[string][]byte{},
		},
	}

	message := "hello, world!!!"
	step := &mockStep{
		name:     "mock_step",
		messages: []string{message},
	}

	steps.RegisterStep(step.Name(), step)

	reqBody := BuildTaskRequest{
		Cfg: steps.Config{
			Timeout: time.Second * 1,
		},
		StepNames: []string{step.Name()},
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

	h.BuildAndRunTask(rec, req)

	resp := &TaskResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), resp)

	if rec.Code != http.StatusCreated {
		t.Errorf("Wrong response code expected %d actual %d", rec.Code, http.StatusCreated)
		return
	}

	if err != nil {
		t.Errorf("Unexpected error while parsing response %v", err)
	}
}
