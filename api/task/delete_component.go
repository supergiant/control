package task

import (
	"encoding/json"

	"github.com/supergiant/supergiant/core"
)

type DeleteComponentMessage struct {
	AppName       string
	ComponentName string
}

// DeleteComponent implements task.Performable interface
type DeleteComponent struct {
	core *core.Core
}

func (j DeleteComponent) Perform(data []byte) error {
	message := new(DeleteComponentMessage)
	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	app, err := j.core.Apps().Get(message.AppName)
	if err != nil {
		return err
	}

	component, err := app.Components().Get(message.ComponentName)
	if err != nil {
		return err
	}

	return component.Delete()
}
