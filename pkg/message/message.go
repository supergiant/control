package message

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
)

// Message is used in response body to display separate messages to user and developers
type Message struct {
	// UserMessage should be used to display message to user on UI/CLI
	UserMessage string `json:"userMessage"`
	// DevMessage should be used to display message to developer on API/console
	DevMessage string `json:"devMessage"`
	// ErrorCode is the unique identifier of an error
	ErrorCode sgerrors.ErrorCode `json:"errorCode"`
	// MoreInfo should be a link to supergiant documentation to display common problems
	MoreInfo string `json:"moreInfo"`
}

func New(userMessage string, devMessage string, code sgerrors.ErrorCode, moreInfo string) Message {
	return Message{
		UserMessage: userMessage,
		DevMessage:  devMessage,
		ErrorCode:   code,
		MoreInfo:    moreInfo,
	}
}
func SendMessage(w http.ResponseWriter, msg Message, status int) {
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func SendInvalidJSON(w http.ResponseWriter, err error) {
	msg := New("User has sent data in malformed format", err.Error(), sgerrors.InvalidJSON, "")
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
}

// SendValidationFailed - this is special case where frontend should parse dev message and present it on the UI
func SendValidationFailed(w http.ResponseWriter, err error) {
	msg := New("Validation Failed", err.Error(), sgerrors.ValidationFailed, "")
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
}

func SendUnknownError(w http.ResponseWriter, err error) {
	msg := New("Internal error occurred, please consult administrator", err.Error(), sgerrors.UnknownError, "")

	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(data)
}

func SendNotFound(w http.ResponseWriter, entityName string, err error) {
	msg := New(fmt.Sprintf("No such %s", entityName), err.Error(), sgerrors.NotFound, "")

	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(data)
}

func SendAlreadyExists(w http.ResponseWriter, entityName string, err error) {
	msg := New(fmt.Sprintf("%s already exists", entityName), err.Error(), sgerrors.AlreadyExists, "")

	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	w.Write(data)
}

func SendInvalidCredentials(w http.ResponseWriter, err error) {
	msg := New("Credentials are bad for cloud provider",
		err.Error(), sgerrors.InvalidCredentials, "")

	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshall message: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
}
