package steps

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/runner"
)

type mockRunner struct {
	mock.Mock
}

func (m *mockRunner) Run(command *runner.Command) error {
	args := m.Called(command)

	val, ok := args.Get(0).(error)

	if !ok {
		return args.Error(0)
	}

	return val
}

func TestRunTemplateSuccess(t *testing.T) {
	r := &mockRunner{}
	r.On("Run", mock.Anything).Return(nil)
	tpl, _ := template.New("noname").Parse("")
	err := RunTemplate(context.Background(), tpl, r,
		&bytes.Buffer{}, &Config{})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestRunTemplateErrRun(t *testing.T) {
	r := &mockRunner{}
	expectedErr := errors.New("message1")
	r.On("Run", mock.Anything).Return(expectedErr)
	tpl, _ := template.New("noname").Parse("")
	actualErr := RunTemplate(context.Background(), tpl, r,
		&bytes.Buffer{}, &Config{})

	if actualErr == nil {
		t.Errorf("Error must not be nil")
	}

	if !strings.Contains(actualErr.Error(), expectedErr.Error()) {
		t.Errorf("Error %v must contain error %v",
			actualErr, expectedErr)
	}
}

func TestRunTemplateDeadline(t *testing.T) {
	r := &mockRunner{}
	r.On("Run", mock.Anything).Return(nil)
	tpl, _ := template.New("noname").Parse("")

	ctx, _ := context.WithDeadline(context.Background(),
		time.Now().Add(-time.Second))
	actualErr := RunTemplate(ctx, tpl, r, &bytes.Buffer{}, &Config{})

	if actualErr == nil {
		t.Errorf("Error must not be nil")
	}

	if !strings.Contains(actualErr.Error(), context.DeadlineExceeded.Error()) {
		t.Errorf("Error %v must contain error %v",
			actualErr, context.DeadlineExceeded.Error())
	}
}

func TestRunTemplateNilWriter(t *testing.T) {
	r := &mockRunner{}
	r.On("Run", mock.Anything).Return(nil)
	tpl, _ := template.New("noname").Parse("")
	actualErr := RunTemplate(context.Background(), tpl, r, nil, &Config{})

	if actualErr == nil {
		t.Errorf("Error must not be nil")
	}

	if !strings.Contains(actualErr.Error(), runner.ErrNilWriter.Error()) {
		t.Errorf("Error %v must contain error %v",
			actualErr, runner.ErrNilWriter.Error())
	}
}

func TestRunTemplateNilContext(t *testing.T) {
	r := &mockRunner{}
	r.On("Run", mock.Anything).Return(nil)
	tpl, _ := template.New("noname").Parse("")
	actualErr := RunTemplate(context.Background(), tpl, r, nil, &Config{})

	if actualErr == nil {
		t.Errorf("Error must not be nil")
	}

	if !strings.Contains(actualErr.Error(), runner.ErrNilWriter.Error()) {
		t.Errorf("Error %v must contain error %v",
			actualErr, runner.ErrNilWriter.Error())
	}
}

func TestRunTemplateErrTemplate(t *testing.T) {
	r := &mockRunner{}
	r.On("Run", mock.Anything).Return(nil)
	tpl, _ := template.New("noname").Parse("{{ .NotFound }}")
	actualErr := RunTemplate(context.Background(), tpl, r, nil, &Config{})

	if actualErr == nil {
		t.Errorf("Error must not be nil")
	}

	if !strings.Contains(actualErr.Error(), "can't evaluate field") {
		t.Errorf("Error %v must contain error %v",
			actualErr, "can't evaluate field")
	}
}
