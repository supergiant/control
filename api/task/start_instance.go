package task

import (
	"encoding/json"

	"github.com/supergiant/supergiant/core"
)

type StartInstanceMessage struct {
	AppName       string
	ComponentName string
	ReleaseID     string
	ID            int
}

// StartInstance implements Performable interface
type StartInstance struct {
	core *core.Core
}

func (j StartInstance) Perform(data []byte) error {
	msg := new(StartInstanceMessage)
	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	app, err := j.core.Apps().Get(msg.AppName)
	if err != nil {
		return err
	}
	component, err := app.Components().Get(msg.ComponentName)
	if err != nil {
		return err
	}
	release, err := component.Releases().Get(msg.ReleaseID)
	if err != nil {
		return err
	}
	instance, err := release.Instances().Get(msg.ID)
	if err != nil {
		return err
	}

	return instance.Start()
}
