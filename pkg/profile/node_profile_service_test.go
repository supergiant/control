package profile

import (
	"context"
	"testing"

	"github.com/pkg/errors"
)

func TestNodeProfileServiceGet(t *testing.T) {
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

	for _, testCase := range testCases {
		storage := &fakeStorage{
			get: func(ctx context.Context, prefix string, key string) ([]byte, error) {
				return testCase.data, testCase.err
			},
		}

		service := NodeProfileService{
			"",
			storage,
		}

		profile, err := service.Get(context.Background(), "fake_id")

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && profile.Id != testCase.expectedId {
			t.Errorf("Wrong profile id expected %s actual %s", testCase.expectedId, profile.Id)
		}
	}
}

func TestNodeProfileServiceCreate(t *testing.T) {
	testCases := []struct {
		node *NodeProfile
		err  error
	}{
		{
			node: &NodeProfile{},
			err:  nil,
		},
		{
			node: &NodeProfile{},
			err:  errors.New("test err"),
		},
	}

	for _, testCase := range testCases {
		storage := &fakeStorage{
			create: func(ctx context.Context, prefix string, key string, value []byte) error {
				return testCase.err
			},
		}

		service := NodeProfileService{
			"",
			storage,
		}

		err := service.Create(context.Background(), testCase.node)

		if testCase.err != err {
			t.Errorf("Unexpected error when create node %v", err)
		}
	}
}

func TestNodeProfileServiceGetAll(t *testing.T) {
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

	for _, testCase := range testCases {
		storage := &fakeStorage{
			getAll: func(ctx context.Context, prefix string) ([][]byte, error) {
				return testCase.data, testCase.err
			},
		}

		service := NodeProfileService{
			"",
			storage,
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
