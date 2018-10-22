package message

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	json "encoding/json"
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestSendUnknownError(t *testing.T) {
	rr := httptest.NewRecorder()
	SendUnknownError(rr, errors.New("test error"))

	msg := new(Message)
	err := jsoniter.NewDecoder(rr.Body).Decode(msg)
	require.NoError(t, err)
	require.Equal(t, sgerrors.UnknownError, msg.ErrorCode)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestNew(t *testing.T) {
	var (
		userMessage = "message"
		devMessage  = "dev_message"
		code        = sgerrors.InvalidJSON
		moreInfo    = "some more information"
	)

	message := New(userMessage, devMessage, code, moreInfo)

	if message.UserMessage != userMessage {
		t.Errorf("Wrong user message expected %s actual %s",
			userMessage, message.UserMessage)
	}

	if message.DevMessage != devMessage {
		t.Errorf("Wonr dev message expected %s actual %s",
			devMessage, message.DevMessage)
	}

	if message.ErrorCode != code {
		t.Errorf("Wrong error code expected %d actual %d",
			code, message.ErrorCode)
	}

	if message.MoreInfo != moreInfo {
		t.Errorf("Wrong additional info expected %s actual %s",
			moreInfo, message.MoreInfo)
	}
}

func TestSendMessage(t *testing.T) {
	expectedCode := http.StatusInternalServerError
	expectedHeaderKey := "Content-Type"
	expectedHeaderValue := "application/json"
	expectedMessage := "userMessage"
	expectedDevMessage := "devMessage"

	rec := httptest.NewRecorder()
	msg := New(expectedMessage, expectedDevMessage,
		sgerrors.InvalidCredentials, "")

	SendMessage(rec, msg, expectedCode)

	if rec.Code != expectedCode {
		t.Errorf("Wrong code expected %d actual %d", expectedCode, rec.Code)
	}

	if h := rec.Header().Get(expectedHeaderKey); h != expectedHeaderValue {
		t.Errorf("Wrong header expected %s actual %s", expectedHeaderValue, h)
	}

	msg2 := &Message{}
	err := json.Unmarshal(rec.Body.Bytes(), msg2)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if msg2.UserMessage != expectedMessage {
		t.Errorf("Wrong user message expected %s actual %s",
			expectedMessage, msg2.UserMessage)
	}

	if msg2.DevMessage != expectedDevMessage {
		t.Errorf("Wonr dev message expected %s actual %s",
			expectedDevMessage, msg2.DevMessage)
	}
}

func TestSendInvalidJSON(t *testing.T) {
	header := "Content-Type"
	headerValue := "application/json"
	errMsg := "expected error dev message"
	err := errors.New(errMsg)
	rec := httptest.NewRecorder()

	SendInvalidJSON(rec, err)

	if h := rec.Header().Get(header); h != headerValue {
		t.Errorf("Wrong header expected %s actual %s", headerValue, h)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Wrong code expected %d actual %d",
			http.StatusBadRequest, rec.Code)
	}

	msg2 := &Message{}
	err = json.Unmarshal(rec.Body.Bytes(), msg2)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if msg2.DevMessage != errMsg {
		t.Errorf("Wonr dev message expected %s actual %s",
			errMsg, msg2.DevMessage)
	}
}

func TestSendValidationFailed(t *testing.T) {
	header := "Content-Type"
	headerValue := "application/json"
	errMsg := "expected error dev message"
	err := errors.New(errMsg)
	rec := httptest.NewRecorder()

	SendValidationFailed(rec, err)

	if h := rec.Header().Get(header); h != headerValue {
		t.Errorf("Wrong header expected %s actual %s", headerValue, h)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Wrong code expected %d actual %d",
			http.StatusBadRequest, rec.Code)
	}

	msg2 := &Message{}
	err = json.Unmarshal(rec.Body.Bytes(), msg2)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if msg2.DevMessage != errMsg {
		t.Errorf("Wonr dev message expected %s actual %s",
			errMsg, msg2.DevMessage)
	}
}

func TestSendNotFound(t *testing.T) {
	header := "Content-Type"
	headerValue := "application/json"
	errMsg := "expected error dev message"
	entityName := "entity"
	err := errors.New(errMsg)
	rec := httptest.NewRecorder()

	SendNotFound(rec, entityName, err)

	if h := rec.Header().Get(header); h != headerValue {
		t.Errorf("Wrong header expected %s actual %s", headerValue, h)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Wrong code expected %d actual %d",
			http.StatusNotFound, rec.Code)
	}

	msg2 := &Message{}
	err = json.Unmarshal(rec.Body.Bytes(), msg2)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if msg2.DevMessage != errMsg {
		t.Errorf("Wonr dev message expected %s actual %s",
			errMsg, msg2.DevMessage)
	}

	if msg2.UserMessage != fmt.Sprintf("No such %s", entityName) {
		t.Errorf("wrong user message expected %s actual %s",
			fmt.Sprintf("No such %s", entityName), msg2.UserMessage)
	}
}
