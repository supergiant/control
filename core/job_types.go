package core

import "encoding/json"

type JobType int

const (
	JobTypeCreateComponent JobType = iota
	JobTypeDestroyComponent
	JobTypeDestroyApp
)

// TODO we can consolidate create/destroy messages for each resource

// Create component Job
//==============================================================================
type CreateComponentMessage struct {
	AppName       string
	ComponentName string
}

// CreateComponent implements job.Performable interface
type CreateComponent struct {
	client *Client
}

func (j CreateComponent) Perform(data []byte) error {
	message := new(CreateComponentMessage)
	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	app, err := j.client.Apps().Get(message.AppName)
	if err != nil {
		return err
	}

	component, err := app.Components().Get(message.ComponentName)
	if err != nil {
		return err
	}

	return component.Provision()
}

// Destroy component Job
//==============================================================================
type DestroyComponentMessage struct {
	AppName       string
	ComponentName string
}

// DestroyComponent implements job.Performable interface
type DestroyComponent struct {
	client *Client
}

func (j DestroyComponent) Perform(data []byte) error {
	message := new(DestroyComponentMessage)
	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	app, err := j.client.Apps().Get(message.AppName)
	if err != nil {
		return err
	}

	component, err := app.Components().Get(message.ComponentName)
	if err != nil {
		return err
	}

	if err := component.Teardown(); err != nil {
		return err
	}

	return component.r.Delete(component.Name)
}

// Destroy app Job
//==============================================================================
type DestroyAppMessage struct {
	AppName string
}

// DestroyApp implements job.Performable interface
type DestroyApp struct {
	client *Client
}

func (j DestroyApp) Perform(data []byte) error {
	message := new(DestroyAppMessage)
	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	app, err := j.client.Apps().Get(message.AppName)
	if err != nil {
		return err
	}

	components, err := app.Components().List()
	if err != nil {
		return err
	}
	for _, component := range components.Items {
		// TODO this should just be a Delete() or Destroy() method (pick one and stick with it) on component
		if err := component.Teardown(); err != nil {
			return err
		}
		if err := component.r.Delete(component.Name); err != nil {
			return err
		}
	}

	if err = app.DeleteNamespace(); err != nil {
		return err
	}

	return app.r.Delete(app.Name)
}
