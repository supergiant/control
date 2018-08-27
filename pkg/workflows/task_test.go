package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type MockRepository struct {
	storage map[string][]byte
}

func (f *MockRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	f.storage[fmt.Sprintf("%s/%s", prefix, key)] = value

	return nil
}

func (f *MockRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	return f.storage[fmt.Sprintf("%s/%s", prefix, key)], nil
}

func (f *MockRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	return nil, nil
}

func (f *MockRepository) Delete(ctx context.Context, prefix string, key string) error {
	return nil
}

type MockStep struct {
	name        string
	description string
	counter     int
	messages    []string
	errs        []error
}

func (f *MockStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	defer func() {
		f.counter++
	}()

	if f.counter < len(f.errs) && f.errs[f.counter] != nil {
		out.Write([]byte(f.errs[f.counter].Error()))
	} else if f.messages != nil && len(f.messages) > f.counter {
		out.Write([]byte(f.messages[f.counter]))
	}

	if f.counter < len(f.errs) {
		return f.errs[f.counter]
	}

	return nil
}

func (f *MockStep) Name() string {
	return f.name
}

func (f *MockStep) Description() string {
	return f.description
}

func (f *MockStep) Depends() []string {
	return nil
}

func TestTaskRunError(t *testing.T) {
	errMsg := "something has gone wrong"
	s := &MockRepository{
		storage: make(map[string][]byte),
	}
	id := "abcd"

	workflow := Task{
		ID:         id,
		repository: s,
		workflow: []steps.Step{
			&MockStep{name: "step1", errs: nil},
			&MockStep{name: "step2", errs: []error{errors.New(errMsg)}},
			&MockStep{name: "step3", errs: nil}},
	}

	buffer := &bytes.Buffer{}
	errChan := workflow.Run(context.Background(), steps.Config{}, buffer)

	if len(workflow.ID) == 0 {
		t.Error("id must not be empty")
	}

	err := <-errChan

	if err == nil {
		t.Error("Error must not be nil")
	}

	if !strings.Contains(buffer.String(), errMsg) {
		t.Error(fmt.Sprintf("Expected error message %s not found in output %s", errMsg, buffer.String()))
	}

	w := &Task{}
	data := s.storage[fmt.Sprintf("%s/%s", prefix, id)]

	err = json.Unmarshal([]byte(data), w)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if w.StepStatuses[1].Status != steps.StatusError {
		t.Errorf("Unexpected step statues expected %s actual %s",
			steps.StatusError, w.StepStatuses[1].Status)
	}
}

func TestTaskRunSuccess(t *testing.T) {
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	id := "abcd"
	task := Task{
		ID:         id,
		repository: s,
		workflow: []steps.Step{
			&MockStep{name: "step1", errs: nil},
			&MockStep{name: "step2", errs: nil},
			&MockStep{name: "step3", errs: nil}},
	}

	buffer := &bytes.Buffer{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)

	if len(id) == 0 {
		t.Error("id must not be empty")
	}

	err := <-errChan

	if err != nil {
		t.Error("Error must be nil")
	}

	w := &Task{}
	data := s.storage[fmt.Sprintf("%s/%s", prefix, task.ID)]

	err = json.Unmarshal([]byte(data), w)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}
	for _, status := range w.StepStatuses {
		if status.Status != steps.StatusSuccess {
			t.Errorf("Unexpected status expectec %s actual %s", steps.StatusSuccess, status.Status)
		}
	}
}

func TestWorkflowRestart(t *testing.T) {
	errMsg := "something has gone wrong"
	s := &MockRepository{
		storage: make(map[string][]byte),
	}
	id := "abcd"

	task := &Task{
		ID:         id,
		repository: s,
		workflow: []steps.Step{
			&MockStep{name: "step1", errs: nil},
			&MockStep{name: "step2", errs: []error{errors.New(errMsg), nil}},
			&MockStep{name: "step3", errs: nil},
		},
	}

	buffer := &bytes.Buffer{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)

	if len(id) == 0 {
		t.Error("id must not be empty")
	}

	err := <-errChan

	if err == nil {
		t.Error("Error must not be nil")
	}

	if !strings.Contains(buffer.String(), errMsg) {
		t.Error(fmt.Sprintf("Expected error message %s not found in output %s", errMsg, buffer.String()))
	}

	data := s.storage[fmt.Sprintf("%s/%s", prefix, id)]
	err = json.Unmarshal([]byte(data), task)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if task.StepStatuses[1].Status != steps.StatusError {
		t.Errorf("Unexpected step statues expected %s actual %s",
			steps.StatusError, task.StepStatuses[1].Status)
	}

	buffer.Reset()
	errChan = task.Restart(context.Background(), id, buffer)
	err = <-errChan

	if err != nil {
		t.Errorf("Error must not be nil actual %v", err)
	}
}
