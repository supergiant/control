package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"context"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
)

func TestWorkflowHandlerGetWorkflow(t *testing.T) {
	id := "abcd"
	expectedType := "MasterTask"
	expectedSteps := []StepStatus{{}, {}}
	w1 := &Task{
		Type:         expectedType,
		StepStatuses: expectedSteps,
	}
	data, _ := json.Marshal(w1)

	h := TaskHandler{
		repository: &MockRepository{
			map[string][]byte{
				fmt.Sprintf("%s/%s", Prefix, id): data,
			},
		},
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", Prefix, id), nil)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s/{id}", Prefix), h.GetTask)
	router.ServeHTTP(resp, req)

	w2 := &Task{}
	err := json.Unmarshal(resp.Body.Bytes(), w2)

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

func TestTaskHandlerRunTask(t *testing.T) {
	Init()
	h := TaskHandler{
		runnerFactory: func(cfg ssh.Config) (runner.Runner, error) {
			return &testutils.MockRunner{}, nil
		},
		repository: &MockRepository{
			map[string][]byte{},
		},
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
	}

	workflowName := "workflow1"
	message := "hello, world!!!"
	step := &MockStep{
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
		t.Errorf("Unexpected err while json encoding req body %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	h.RunTask(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Wrong response code expected %d received %d", http.StatusAccepted, rec.Code)
	}

	resp := &TaskResponse{}
	json.NewDecoder(rec.Body).Decode(resp)

	if len(resp.ID) == 0 {
		t.Error("task id in response should not be empty")
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
			SshConfig: steps.SshConfig{
				User: "root",
				Port: "22",
				PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
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

func TestWorkflowHandlerBuildWorkflow(t *testing.T) {
	h := TaskHandler{
		runnerFactory: func(cfg ssh.Config) (runner.Runner, error) {
			return &testutils.MockRunner{}, nil
		},
		repository: &MockRepository{
			map[string][]byte{},
		},
	}

	message := "hello, world!!!"
	step := &MockStep{
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
		t.Errorf("Unexpected err while parsing response %v", err)
	}
}
