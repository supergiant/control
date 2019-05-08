package templatemanager

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/templates"
)

var (
	m           sync.RWMutex
	templateMap map[string]*template.Template
)

func init() {
	templateMap = make(map[string]*template.Template)
}

func Init(templateDir string) error {
	if err := addDefaultTpls(); err != nil {
		return errors.Wrap(err, "add default templates")
	}
	if err := addCustomTpls(templateDir); err != nil {
		return errors.Wrap(err, "add custom templates")
	}
	return nil
}

func GetTemplate(templateName string) (*template.Template, error) {
	m.RLock()
	m.RUnlock()
	if tpl := templateMap[templateName]; tpl != nil {
		return tpl, nil
	} else {
		return nil, sgerrors.ErrNotFound
	}
}

func SetTemplate(templateName string, tpl *template.Template) {
	m.Lock()
	defer m.Unlock()
	templateMap[templateName] = tpl
}

func DeleteTemplate(templateName string) {
	m.Lock()
	defer m.Unlock()
	delete(templateMap, templateName)
}

func addDefaultTpls() error {
	for name, tpl := range templates.Default {
		t, err := template.New(name).Parse(tpl)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s template", name)
		}
		logrus.Debugf("templatemanager: adding default template: %q", name)
		SetTemplate(name, t)
	}
	return nil
}

func addCustomTpls(dirname string) error {
	if len(dirname) == 0 {
		return nil
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			fullName := path.Join(dirname, f.Name())

			f, err := os.Open(fullName)

			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(f)

			if err != nil {
				return err
			}

			lastTerm := len(strings.Split(f.Name(), "/"))
			name := strings.Split(strings.Split(f.Name(), "/")[lastTerm-1], ".")[0]

			t, err := template.New(name).Parse(string(data))
			if err != nil {
				return errors.Wrapf(err, "failed to parse %s template", fullName)
			}

			logrus.Debugf("templatemanager: adding custom template: %q", name)
			templateMap[name] = t
		}
	}

	return nil
}
