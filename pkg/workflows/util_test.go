package workflows

import (
	"strings"
	"testing"

	"github.com/supergiant/control/pkg/util"
)

func TestBindParams(t *testing.T) {
	obj := &struct {
		ParamA string `json:"a"`
		ParamB string `json:"b"`
	}{}

	expectedA := "hello"
	expectedB := "world"

	params := map[string]string{
		"a": expectedA,
		"b": expectedB,
	}

	util.BindParams(params, obj)

	if !strings.EqualFold(obj.ParamA, params["a"]) {
		t.Errorf("Wrong value for paramA expected %s actual %s", params["a"], obj.ParamA)
	}

	if !strings.EqualFold(obj.ParamB, params["b"]) {
		t.Errorf("Wrong value for paramB expected %s actual %s", params["b"], obj.ParamB)
	}
}
