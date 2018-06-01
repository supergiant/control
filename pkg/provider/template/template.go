package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/provider"
)

type templateMap struct {
	m map[string]*template.Template
}

var (
	Templates           templateMap
	errTemplateNotFound = errors.New("template not found")
)

func init() {
	kubeVersions := provider.GetK8sVersions()
	NodeRoles := provider.GetNodeRoles()
	Templates = templateMap{make(map[string]*template.Template)}

	for _, nodeRole := range NodeRoles {
		for _, kubeVersion := range kubeVersions {
			mversion := strings.Split(string(kubeVersion), ".")
			fileName := fmt.Sprintf("config/providers/common/%s.%s/%s.yaml)", mversion[0], mversion[1], string(nodeRole))

			// Create minion template
			userdataTemplate, err := bindata.Asset(fileName)
			if err != nil {
				logrus.Fatalf("Error binding data")
			}
			tpl, err := template.New(fmt.Sprintf("%s_template", string(nodeRole))).Parse(string(userdataTemplate))

			if err != nil {
				logrus.Fatalf("Error creating %s template for %s", string(nodeRole), kubeVersions)
			}

			Templates.m[fileName] = tpl
		}
		// GCE case create either master or minion
		fileName := fmt.Sprintf("config/providers/gce/%s.yaml)", string(nodeRole))
		userdataTemplate, err := bindata.Asset(fileName)
		if err != nil {
			logrus.Fatalf("Error binding data")
		}
		tpl, err := template.New(fmt.Sprintf("%s_template", string(nodeRole))).Parse(string(userdataTemplate))

		if err != nil {
			logrus.Fatalf("Error creating %s template for GCE %s", string(nodeRole))
		}

		Templates.m[fileName] = tpl
	}
}

func (t templateMap) Get(templateName string) (*template.Template, error) {
	tpl, ok := t.m[templateName]

	if !ok {
		return nil, errTemplateNotFound
	}

	return tpl, nil
}
