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

type fakeRepository struct {
	storage map[string][]byte
}

func (f *fakeRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	f.storage[fmt.Sprintf("%s/%s", prefix, key)] = value

	return nil
}

func (f *fakeRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	return f.storage[fmt.Sprintf("%s/%s", prefix, key)], nil
}

func (f *fakeRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	return nil, nil
}

func (f *fakeRepository) Delete(ctx context.Context, prefix string, key string) error {
	return nil
}

type fakeStep struct {
	name        string
	description string
	counter     int
	errs        []error
}

func (f *fakeStep) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	defer func() {
		f.counter++
	}()

	if f.counter < len(f.errs) && f.errs[f.counter] != nil {
		out.Write([]byte(f.errs[f.counter].Error()))
	}

	if f.counter < len(f.errs) {
		return f.errs[f.counter]
	}

	return nil
}

func (f *fakeStep) Name() string {
	return f.name
}

func (f *fakeStep) Description() string {
	return f.description
}

func TestTaskRunError(t *testing.T) {
	errMsg := "something has gone wrong"
	s := &fakeRepository{
		storage: make(map[string][]byte),
	}
	id := "abcd"

	workflow := Task{
		Id:         id,
		Config:     steps.Config{},
		repository: s,
		workflow: []steps.Step{
			&fakeStep{name: "step1", errs: nil},
			&fakeStep{name: "step2", errs: []error{errors.New(errMsg)}},
			&fakeStep{name: "step3", errs: nil}},
	}

	buffer := &bytes.Buffer{}
	errChan := workflow.Run(context.Background(), buffer)

	if len(workflow.Id) == 0 {
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
	s := &fakeRepository{
		storage: make(map[string][]byte),
	}

	id := "abcd"
	workflow := Task{
		Id:         id,
		Config:     steps.Config{},
		repository: s,
		workflow: []steps.Step{
			&fakeStep{name: "step1", errs: nil},
			&fakeStep{name: "step2", errs: nil},
			&fakeStep{name: "step3", errs: nil}},
	}

	buffer := &bytes.Buffer{}
	errChan := workflow.Run(context.Background(), buffer)

	if len(id) == 0 {
		t.Error("id must not be empty")
	}

	err := <-errChan

	if err != nil {
		t.Error("Error must be nil")
	}

	w := &Task{}
	data := s.storage[fmt.Sprintf("%s/%s", prefix, workflow.Id)]

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
	s := &fakeRepository{
		storage: make(map[string][]byte),
	}
	id := "abcd"

	w := &Task{
		Id:         id,
		Config:     steps.Config{},
		repository: s,
		workflow: []steps.Step{
			&fakeStep{name: "step1", errs: nil},
			&fakeStep{name: "step2", errs: []error{errors.New(errMsg), nil}},
			&fakeStep{name: "step3", errs: nil},
		},
	}

	buffer := &bytes.Buffer{}
	errChan := w.Run(context.Background(), buffer)

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
	err = json.Unmarshal([]byte(data), w)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if w.StepStatuses[1].Status != steps.StatusError {
		t.Errorf("Unexpected step statues expected %s actual %s",
			steps.StatusError, w.StepStatuses[1].Status)
	}

	buffer.Reset()
	errChan = w.Restart(context.Background(), id, buffer)
	err = <-errChan

	if err != nil {
		t.Errorf("Error must not be nil actual %v", err)
	}
}
