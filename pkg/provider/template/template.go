package template

import (
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/bindata"
)

var (
	Templates map[string]*template.Template

	kubeVersions []string
)

func init() {
	// TODO(stgleb): We need single place to have all supported k8s versions
	kubeVersions = []string{"1.5", "1.6", "1.7", "1.8"}
	Templates = make(map[string]*template.Template)

	for _, kubeVersion := range kubeVersions {
		mversion := strings.Split(kubeVersion, ".")
		userdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/minion.yaml")
		if err != nil {
			logrus.Fatalf("Error binding data")
		}
		tpl, err := template.New("minion_template").Parse(string(userdataTemplate))

		if err != nil {
			logrus.Fatalf("Error creating template for %s", kubeVersions)
		}

		Templates[kubeVersion] = tpl
	}
}
