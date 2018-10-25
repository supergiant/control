package node

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/testutils"
)

func TestNewService(t *testing.T) {
	repo := &testutils.MockStorage{}
	prefix := "prefix"

	s := NewService(prefix, repo)

	if s.prefix != prefix {
		t.Errorf("expected prefix %s actuall %s", prefix, s.prefix)
	}

	if s.repository != repo {
		t.Errorf("expected repo %v actual %v", repo, s.repository)
	}
}

func TestNodeServiceGet(t *testing.T) {
	testCases := []struct {
		expectedId string
		data       []byte
		err        error
	}{
		{
			expectedId: "1234",
			data:       []byte(`{"id": "1234", "created_at": 1234, "provider": "aws", "region": "us-west1"}`),
			err:        nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := "/node/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), prefix, "fake_id").Return(testCase.data, testCase.err)

		service := Service{
			prefix,
			m,
		}

		node, err := service.Get(context.Background(), "fake_id")

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && node.ID != testCase.expectedId {
			t.Errorf("Wrong node id expected %s actual %s", testCase.expectedId, node.ID)
		}
	}
}

func TestNodeCreate(t *testing.T) {
	testCases := []struct {
		node *Node
		err  error
	}{
		{
			node: &Node{},
			err:  nil,
		},
		{
			node: &Node{},
			err:  errors.New("test err"),
		},
	}

	prefix := "/node/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		kubeData, _ := json.Marshal(testCase.node)

		m.On("Put",
			context.Background(),
			prefix,
			testCase.node.ID,
			kubeData).
			Return(testCase.err)

		service := Service{
			prefix,
			m,
		}

		err := service.Create(context.Background(), testCase.node)

		if testCase.err != err {
			t.Errorf("Unexpected error when create node %v", err)
		}
	}
}

func TestNodeListAll(t *testing.T) {
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

	prefix := "/node/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("GetAll", context.Background(), prefix).Return(testCase.data, testCase.err)

		service := Service{
			prefix,
			m,
		}

		nodes, err := service.ListAll(context.Background())

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && len(nodes) != 2 {
			t.Errorf("Wrong len of nodes expected 2 actual %d", len(nodes))
		}
	}
}

func TestService_Delete(t *testing.T) {
	testCases := []struct {
		nodeId      string
		expectedErr error
	}{
		{
			nodeId:      "node-id",
			expectedErr: errors.New("test"),
		},
		{
			nodeId:      "node-id",
			expectedErr: nil,
		},
	}

	prefix := "prefix"

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}
		mockRepo.On("Delete", mock.Anything,
			prefix, mock.Anything).
			Return(testCase.expectedErr)

		svc := Service{
			repository: mockRepo,
			prefix:     prefix,
		}

		svc.Delete(context.Background(),
			testCase.nodeId)
	}
}
