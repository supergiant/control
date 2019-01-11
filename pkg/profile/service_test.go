package profile

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/testutils"
)

func TestKubeProfileServiceGet(t *testing.T) {
	testCases := []struct {
		expectedId string
		data       []byte
		err        error
	}{
		{
			expectedId: "1234",
			data:       []byte(`{"id":"1234", "nodes":[{},{}]}`),
			err:        nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := "/profile/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), prefix, "fake_id").Return(testCase.data, testCase.err)

		service := Service{
			prefix,
			m,
		}

		profile, err := service.Get(context.Background(), "fake_id")

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && profile.ID != testCase.expectedId {
			t.Errorf("Wrong profile id expected %s actual %s", testCase.expectedId, profile.ID)
		}
	}
}

func TestKubeProfileServiceCreate(t *testing.T) {
	testCases := []struct {
		profile *Profile
		err     error
	}{
		{
			profile: &Profile{},
			err:     nil,
		},
		{
			profile: &Profile{},
			err:     errors.New("test err"),
		},
	}

	prefix := "/profile/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		kubeData, _ := json.Marshal(testCase.profile)

		m.On("Put",
			context.Background(),
			prefix,
			mock.Anything,
			kubeData).
			Return(testCase.err)

		service := Service{
			prefix,
			m,
		}

		err := service.Create(context.Background(), testCase.profile)

		if testCase.err != err {
			t.Errorf("Unexpected error when create node %v", err)
		}
	}
}

func TestKubeProfileServiceGetAll(t *testing.T) {
	testCases := []struct {
		data [][]byte
		err  error
	}{
		{
			data: [][]byte{[]byte(`{"id":"1234", "nodes":[{},{}]}`), []byte(`{"id":"5678", "nodes":[{},{}]}`)},
			err:  nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := "/profile/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("GetAll", context.Background(), prefix).Return(testCase.data, testCase.err)

		service := Service{
			prefix,
			m,
		}

		profiles, err := service.GetAll(context.Background())

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && len(profiles) != 2 {
			t.Errorf("Wrong len of profiles expected 2 actual %d", len(profiles))
		}
	}
}

func TestNewKubeProfileService(t *testing.T) {
	prefix := "prefix"
	repo := &testutils.MockStorage{}

	svc := NewService(prefix, repo)

	if svc == nil {
		t.Error("service must not be nil")
	}

	if svc.prefix != prefix {
		t.Errorf("Wrong prefix expected %s actual %s", prefix, svc.prefix)
	}

	if svc.kubeProfileStorage != repo {
		t.Errorf("Wrong repo expected %v actual %v", repo, svc.kubeProfileStorage)
	}
}
