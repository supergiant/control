package templatemanager

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

var (
	m sync.RWMutex
	templateMap map[string]*template.Template
)

func init() {
	templateMap = make(map[string]*template.Template)
}

func Init(templateDir string) error {
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			fullName := path.Join(templateDir, f.Name())

			f, err := os.Open(fullName)

			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(f)

			if err != nil {
				return err
			}

			lastTerm := len(strings.Split(f.Name(), "/"))
			key := strings.Split(strings.Split(f.Name(), "/")[lastTerm-1], ".")[0]
			t, _ := template.New(key).Parse(string(data))

			if err != nil {
				return err
			}

			templateMap[key] = t
		}
	}

	return nil
}

func GetTemplate(templateName string) (*template.Template, error) {
	m.RLock()
	m.RUnlock()
	if tpl, ok := templateMap[templateName]; ok {
		return tpl, nil
	} else {
		return nil, sgerrors.ErrNotFound
	}
}

func SetTemplate(templateName string, tpl *template.Template)  {
	m.Lock()
	defer m.Unlock()
	templateMap[templateName] = tpl
}

func DeleteTemplate(templateName string) {
	m.Lock()
	defer m.Unlock()
	delete(templateMap, templateName)
}