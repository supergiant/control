package kube

import (
	"github.com/pkg/errors"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestRestClientForGroupVersion(t *testing.T) {
	testCases := []struct {
		kube *model.Kube
		gv   schema.GroupVersion

		expectedErr error
	}{
		{
			kube:        &model.Kube{},
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			kube: &model.Kube{
				Masters: map[string]*node.Node{
					"node-1": {
						Name:     "node-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		client, err := restClientForGroupVersion(testCase.kube, testCase.gv)

		if errors.Cause(err) != testCase.expectedErr {
			t.Errorf("expected error %v actual %v",
				testCase.expectedErr, err)
		}

		if testCase.expectedErr == nil && client == nil {
			t.Error("client must not be nil")
		}
	}
}

func TestDiscoveryClient(t *testing.T) {
	testCases := []struct {
		kube        *model.Kube
		expectedErr error
	}{
		{
			kube:        &model.Kube{},
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			kube: &model.Kube{
				Masters: map[string]*node.Node{
					"node-1": {
						Name:     "node-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		client, err := discoveryClient(testCase.kube)

		if errors.Cause(err) != testCase.expectedErr {
			t.Errorf("expected error %v actual %v",
				testCase.expectedErr, err)
		}

		if testCase.expectedErr == nil && client == nil {
			t.Error("client must not be nil")
		}
	}
}
