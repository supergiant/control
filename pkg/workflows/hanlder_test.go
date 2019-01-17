package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hpcloud/tail"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/testutils"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockCloudAccountService struct {
	cloudAccount *model.CloudAccount
	err          error
}

func (m *mockCloudAccountService) Get(ctx context.Context, name string) (*model.CloudAccount, error) {
	return m.cloudAccount, m.err
}

func TestWorkflowHandlerGetWorkflow(t *testing.T) {
	id := "abcd"
	expectedType := "MasterTask"
	expectedSteps := []StepStatus{{}, {}}
	w1 := &Task{
		Type:         expectedType,
		StepStatuses: expectedSteps,
	}
	data, err := json.Marshal(w1)

	if err != nil {
		t.Errorf("json marshall %v", err)
	}

	h := TaskHandler{
		repository: &MockRepository{
			map[string][]byte{
				fmt.Sprintf("%s%s", Prefix, id): data,
			},
		},
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", Prefix, id), nil)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s/{id}", Prefix), h.GetTask)
	router.ServeHTTP(resp, req)

	w2 := &Task{}
	err = json.Unmarshal(resp.Body.Bytes(), w2)

	if err != nil {
		t.Errorf("Unexpected err %v", err)
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

func TestTaskHandlerRestartTask(t *testing.T) {
	Init()
	repository := &MockRepository{
		make(map[string][]byte),
	}
	h := TaskHandler{
		runnerFactory: func(cfg ssh.Config) (runner.Runner, error) {
			return &testutils.MockRunner{}, nil
		},
		repository: repository,
		cloudAccGetter: &mockCloudAccountService{
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"fingerprints": "fingerprint",
					"accessToken":  "abcd",
				},
			},
			err: nil,
		},
		getWriter: func(id string) (io.WriteCloser, error) {
			return &bufferCloser{}, nil
		},
	}

	workflowName := "workflow1"
	message := "hello, world!!!"
	step := &MockStep{
		name:     "mock_step",
		messages: []string{message},
	}

	workflow := []steps.Step{step}
	RegisterWorkFlow(workflowName, workflow)

	wf := GetWorkflow(workflowName)

	taskId := "1234"

	task := &Task{
		ID: taskId,
		Config: &steps.Config{
			Kube: model.Kube{
				SSHConfig: model.SSHConfig{
					User: "root",
					Port: "22",
					BootstrapPrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAtArxGzmUffkRNy4bpITg0oicUA6itrh2RumMoydra2QqRL8i
sA6xBaPHbBAOJO/gY/h/qvr8Hnb38GFJcQQy2eENb83i2u8BVnnN2IFkgyCyYCN7
DE54bQejH0xD4qMhXdyEUOyKaOBzHHBliyIR4HmobiddJho4G0Ku3onLDm+++XNG
ZVNj0drFOE0YG+s/Zy5j/8EH3b2NgNzoE9h+jcIkZaRNuYYo5e26RCaakXeJT/Iu
vRFYtSPy3tgGQ+Q/Aj1Kv6Gjv0OFTqO7mmHN+nShThtjdLaCMM8hr72OoVkjljk1
kUBG7VAy0YlgRzxlEOUyXGKY7cOU4Nx51RIYQwIDAQABAoIBAEV2d0F+vKjBoH++
nVGjJq5zoINOsj52+sMvNmB4Q/yB/8DYUYTFlkzLvJQXua1MkzFe3brU7NLAKban
glNFQG1JZAq/z4eScNyxT9b5TRM+WTO4XLAJ0nKWYLwhi4t0TtpMywwBxwDhn+fY
AYVllqoZpf8h1tFtijoSRy960En381UuXWMGUDwzpr01G1GLX0ux4sCVMB4l1SX2
yWxtS/tJKvdLUNlEH72n+w9uz4Xpt62QXubDiibFDtomDbJrIi15umMAGcSgIvLq
2T/HoyXysIKdJNzu8jTYLBHY0oKmMN+DS047rDL6HMX3OH2GTpbnWBt8vk9Ndgwh
akMvfzkCgYEA35icpBVRDS+NUPMFbWcZRPkv8mSCPxJBL2/PHhRyjf85lOiMG7Rl
FVNA1pvATkD5QRo1SP4QiMcLbtLxX5bcevmg48ziGoWcSync5PUKnSwKMR/AAolZ
HD8pL2sjXpedNy1twwIY4XWALC+jf3h2UlBD0Xc/WcHPvpjNaPH1QiUCgYEAziKB
IuagG2uUkeBW5iQiC4eGuelTnRgCBN8CO6TWL88q0IEHtgx9c/s+rQ1AeR5VbGRu
lECDKCw6aq6swAsW4sacGYxvMyivdnSKSrKTFqDbV8ccJ+g30HoVTrcfXN1pZIpd
urEFVTi7uokINLuhVHqmY9poOkLD9ugx1G1dwEcCgYEAjWLyQetcyiq0gGh7mRdl
ajDr+alGlt1TLMzVuh6R5WprHdcCqY4jkR2I1Wu9aX46XslUwmgtSmAawaRPjvNV
TcnFy+ZFXyH3l6vMC1dLs+EiPLfn8XKqT2s8/sgPoIPcnQRz8KjF1OM4/jfNehBR
OXGZwL1X6MjeNZZn5SGCxaUCgYEAwpkGShloeppT9mbQApGH5lR6FpYzzjD07v5M
0FurrBSqOY4l5nHiGRTNtXa2L+E4CYzCa4h/iPQ/7aibAu01HL8cbG4MKEK3al9Y
km4Et68BgttANFhgIJqv9NChdfy72yNYmr805qAZcV6d9ZJQGj1zSP7NuHqBH11S
dVUN1U8CgYEAw5N6ScysYb9Jsaurcykij4mn1tvXzpDcap/Lqu/QXSUJZU1D7Cac
OOJSve1MuYQbV1LEIc15yMPsWTTik2Z98r9IL+3xdofh9yFaG1nxzi9OkN6aVMAz
dZM6MSCYh9kcT0pi2FPmY9iXba9kx4XAnf+0YB5xCz9QSMk4W5xSTBs=
-----END RSA PRIVATE KEY-----`,
					Timeout: 10,
				},
			},
			Node: node.Node{
				PublicIp: "10.20.30.40",
			},
		},

		workflow:   wf,
		repository: repository,
	}

	data, _ := json.Marshal(task)
	repository.Put(context.Background(), Prefix, task.ID, data)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s/restart", Prefix, taskId), nil)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s/{id}/restart", Prefix), h.RestartTask)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Wrong response code expected %d received %d", http.StatusAccepted, rec.Code)
	}
}

func TestTaskHandler_GetLogs(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/abcd/logs/ws", nil)

	rec.Header().Add("Authorization", "bearer none")
	router := mux.NewRouter()
	handler := TaskHandler{}
	handler.Register(router)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Wrong status code expected %d actual %d",
			http.StatusOK, rec.Code)
	}
}

func TestTaskHandler_StreamLogs(t *testing.T) {
	testCases := []struct {
		description  string
		getTailErr   error
		t            *tail.Tail
		expectedCode int
	}{
		{
			description:  "wrong task id",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "wrong task id",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "file not found",
			getTailErr:   os.ErrNotExist,
			expectedCode: http.StatusNotFound,
		},
		{
			description:  "unknown error",
			getTailErr:   errors.New("error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			description: "upgrade error",
			t: &tail.Tail{
				Lines: make(chan *tail.Line),
			},

			expectedCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/tasks/abcd/logs", nil)

		router := mux.NewRouter()
		handler := TaskHandler{
			getTail: func(s string) (*tail.Tail, error) {
				return testCase.t, testCase.getTailErr
			},
		}
		handler.Register(router)
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestNewTaskHandler(t *testing.T) {
	r := &testutils.MockStorage{}
	h := NewTaskHandler(r, nil, nil)

	if h == nil {
		t.Errorf("Handler must not be nil")
	}
}
