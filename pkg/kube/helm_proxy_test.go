package kube

import (
	"strings"
	"testing"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/sgerrors"
)

func TestHelmProxyFrom(t *testing.T) {
	testCases := []struct {
		description string
		k           *model.Kube
		errMsg      string
	}{
		{
			description: "nil kube",
			errMsg:      sgerrors.ErrNilEntity.Error(),
		},
		{
			description: "no master",
			k: &model.Kube{
				Masters: map[string]*node.Node{},
			},
			errMsg: sgerrors.ErrNotFound.Error(),
		},
		{
			description: "success",
			k: &model.Kube{
				Masters: map[string]*node.Node{
					"key": {
						ID:       "key",
						Name:     "key",
						PublicIp: "10.20.30.40",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		_, err := helmProxyFrom(testCase.k)

		if err == nil && testCase.errMsg != "" {
			t.Error("err must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error %v must contain %s", err, testCase.errMsg)
		}
	}
}
