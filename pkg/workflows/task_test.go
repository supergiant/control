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

	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type bufferCloser struct {
	bytes.Buffer
	err error
}

func (b *bufferCloser) Close() error {
	return b.err
}

type MockRepository struct {
	storage map[string][]byte
}

func (f *MockRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	f.storage[prefix+key] = value

	return nil
}

func (f *MockRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	return f.storage[prefix+key], nil
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
	rollback    bool
}

func (f *MockStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	f.rollback = true
	return nil
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

func TestNewTask(t *testing.T) {
	mockRepository := &MockRepository{
		storage: map[string][]byte{},
	}

	workflowMap = make(map[string]Workflow)
	RegisterWorkFlow(ProvisionMaster, Workflow{})

	testCases := []struct {
		taskType      string
		expectedError error
	}{
		{
			ProvisionMaster,
			nil,
		},
		{
			"foo",
			sgerrors.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		task, err := NewTask(testCase.taskType, mockRepository)

		if err != testCase.expectedError {
			t.Errorf("Unexpected error %v", err)
			return
		}

		if testCase.expectedError == nil {
			if _, ok := mockRepository.storage[Prefix+task.ID]; !ok {
				t.Errorf("Task %s found in storage", task.ID)
			}
		}
	}
}

func TestTaskRunError(t *testing.T) {
	errMsg := "something has gone wrong"
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	wf := []steps.Step{
		&MockStep{name: "step1", errs: nil},
		&MockStep{name: "step2", errs: []error{errors.New(errMsg)}},
		&MockStep{name: "step3", errs: nil},
	}

	workflowMap = make(map[string]Workflow)
	RegisterWorkFlow("mock", wf)
	task, err := NewTask("mock", s)

	buffer := &bufferCloser{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)

	err = <-errChan

	if err == nil {
		t.Error("Error must not be nil")
	}

	if !strings.Contains(buffer.String(), errMsg) {
		t.Error(fmt.Sprintf("Expected error message %s not found in output %s", errMsg, buffer.String()))
	}

	w := &Task{}
	data := s.storage[Prefix+task.ID]

	err = json.Unmarshal([]byte(data), w)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if w.Status != statuses.Error {
		t.Errorf("Unexpected task status expected %s actual %s",
			statuses.Error, w.Status)
	}

	if w.StepStatuses[1].Status != statuses.Error {
		t.Errorf("Unexpected step statues expected %s actual %s",
			statuses.Error, w.StepStatuses[1].Status)
	}
}

func TestTaskRunSuccess(t *testing.T) {
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	wf := []steps.Step{
		&MockStep{name: "step1", errs: nil},
		&MockStep{name: "step2", errs: nil},
		&MockStep{name: "step3", errs: nil},
	}

	workflowMap = make(map[string]Workflow)
	RegisterWorkFlow("mock", wf)
	task, err := NewTask("mock", s)

	buffer := &bufferCloser{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)

	err = <-errChan

	if err != nil {
		t.Error("Error must be nil")
	}

	w := &Task{}
	data := s.storage[Prefix+task.ID]

	err = json.Unmarshal([]byte(data), w)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if w.Status != statuses.Success {
		t.Errorf("Unexpected task status expected %s actual %s", statuses.Success, w.Status)
	}
	for _, status := range w.StepStatuses {
		if status.Status != statuses.Success {
			t.Errorf("Unexpected status expectec %s actual %s", statuses.Success, status.Status)
		}
	}
}

func TestWorkflowRestart(t *testing.T) {
	errMsg := "something has gone wrong"
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	wf := []steps.Step{
		&MockStep{name: "step1", errs: nil},
		&MockStep{name: "step2", errs: []error{errors.New(errMsg), nil}},
		&MockStep{name: "step3", errs: nil},
	}
	workflowMap = make(map[string]Workflow)
	RegisterWorkFlow("mock", wf)
	task, err := NewTask("mock", s)

	if err != nil {
		t.Errorf("Error while creating task %v", err)
		return
	}

	buffer := &bufferCloser{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)
	err = <-errChan

	if err == nil {
		t.Error("Error must not be nil")
	}

	if !strings.Contains(buffer.String(), errMsg) {
		t.Error(fmt.Sprintf("Expected error message %s not found in output %s", errMsg, buffer.String()))
	}

	data := s.storage[Prefix+task.ID]
	err = json.Unmarshal([]byte(data), task)

	if err != nil {
		t.Errorf("Unexpected error while unmarshalling data %v", err)
	}

	if task.StepStatuses[1].Status != statuses.Error {
		t.Errorf("Unexpected step statues expected %s actual %s",
			statuses.Error, task.StepStatuses[1].Status)
	}

	buffer.Reset()
	errChan = task.Run(context.Background(), steps.Config{}, buffer)
	err = <-errChan

	if err != nil {
		t.Errorf("Error must not be nil actual %v", err)
	}
}

func TestRollback(t *testing.T) {
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	mockStep := &MockStep{name: "step1", errs: []error{errors.New("should happen")}}
	wf := []steps.Step{
		mockStep,
	}
	workflowMap = make(map[string]Workflow)
	RegisterWorkFlow("mock", wf)
	task, err := NewTask("mock", s)

	require.False(t, mockStep.rollback)

	buffer := &bufferCloser{}
	errChan := task.Run(context.Background(), steps.Config{}, buffer)
	err = <-errChan
	require.Error(t, err)

	require.True(t, mockStep.rollback)
}

type PanicStep struct {
}

func (PanicStep) Run(context.Context, io.Writer, *steps.Config) error {
	panic("implement me")
}

func (PanicStep) Name() string {
	panic("implement me")
}

func (PanicStep) Description() string {
	panic("implement me")
}

func (PanicStep) Depends() []string {
	panic("implement me")
}

func (PanicStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	panic("implement me")
}

func TestPanicHandler(t *testing.T) {
	s := &MockRepository{
		storage: make(map[string][]byte),
	}

	step := &PanicStep{}
	task := &Task{
		ID:         "XYZ",
		repository: s,
		workflow: []steps.Step{
			step,
		},
	}
	errChan := task.Run(context.Background(), steps.Config{}, &bufferCloser{})

	err := <-errChan
	require.Error(t, err)
}
