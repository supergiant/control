package template

import (
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	"fmt"

	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/provider"
)

var (
	Templates map[string]*template.Template
)

func init() {
	kubeVersions := []provider.K8SVersion{provider.K8S15, provider.K8S16, provider.K8S17, provider.K8S18}
	NodeRoles := []provider.NodeRole{provider.Master, provider.Minion}
	Templates = make(map[string]*template.Template)

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

			Templates[fileName] = tpl
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

		Templates[fileName] = tpl
	}
}
