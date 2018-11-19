package templatemanager

import (
	"testing"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestGetTemplateNotFound(t *testing.T) {
	_, err := GetTemplate("not_found.sh.tpl")

	if !sgerrors.IsNotFound(err) {
		t.Errorf("Wrong error expected %v actual %v", sgerrors.ErrNotFound, err)
	}
}
