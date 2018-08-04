package templatemanager

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	templateMap map[string]*template.Template
)

func Init(templateDir string) error {
	templateMap = make(map[string]*template.Template)
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
			key := strings.Split(strings.Split(f.Name(), "/")[lastTerm - 1], ".")[0]
			t, _ := template.New(key).Parse(string(data))

			if err != nil {
				return err
			}

			templateMap[key] = t
		}
	}

	return nil
}

func GetTemplate(templateName string) *template.Template {
	return templateMap[templateName]
}
