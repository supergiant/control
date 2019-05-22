package templatemanager

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/templates"
)

func TestGetTemplate(t *testing.T) {
	testKey := "testGetKey"
	testValue := &template.Template{}
	templateMap[testKey] = testValue

	_, err := GetTemplate("testGetKey")

	if sgerrors.IsNotFound(err) {
		t.Errorf("Template %s not found", testKey)
	}
}

func TestGetTemplateNotFound(t *testing.T) {
	_, err := GetTemplate("not_found.sh.tpl")

	if !sgerrors.IsNotFound(err) {
		t.Errorf("Wrong error expected %v actual %v", sgerrors.ErrNotFound, err)
	}
}

func TestDeleteTemplate(t *testing.T) {
	testKey := "testDeleteKey"
	testValue := &template.Template{}
	templateMap[testKey] = testValue

	DeleteTemplate(testKey)

	if _, ok := templateMap[testKey]; ok {
		t.Errorf("key %s must not be in templateMap %v", testKey, templateMap)
	}
}

func TestSetTemplate(t *testing.T) {
	testKey := "testSetKey"
	testValue := &template.Template{}
	templateMap[testKey] = testValue

	SetTemplate(testKey, testValue)

	if _, ok := templateMap[testKey]; !ok {
		t.Errorf("key %s must not in templateMap %v", testKey, templateMap)
	}
}

func TestInitDefault(t *testing.T) {
	err := Init("")
	require.Nil(t, err)

	for k := range templates.Default {
		_, err := GetTemplate(k)
		require.Nilf(t, err, "%q template should exist")
	}
}
