package jobs

import (
	"io/ioutil"
	"text/template"
)

func ReadTemplate(fileName, templateName string) (*template.Template, error) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	tpl, err := template.New(templateName).Parse(string(data))

	if err != nil {
		return nil, err
	}

	return tpl, nil
}
