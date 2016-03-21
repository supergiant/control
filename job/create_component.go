package job

import (
	"encoding/json"
	"supergiant/core/model"
)

type CreateComponentMessage struct {
	AppName       string
	ComponentName string
}

// CreateComponent implements job.Performable interface
type CreateComponent struct {
	client *model.Client // TODO resources would be better than model semantically

	// if you store variables here, you can define each bit of work as an anonymous function of the job, giving you progress output

	// ------------------- and then, the interface should not expose Perform(), but instead should provide a Performable instance,
	// 																																									that responds to Actions(), returning an array of Action (interface),
	// 																																									that must respond to Perform()

	// ---------------------------------------------------------------------------    actions can be defined with an ? arbitrary string ? that denotes concurrent groups

}

func (j CreateComponent) MaxAttempts() int {
	return 20
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

	if err := component.ProvisionSecrets(); err != nil {
		return err
	}

	// Create Services
	if component.HasInternalPorts() {
		if err := component.ProvisionInternalService(); err != nil {
			return err
		}
	}
	if component.HasExternalPorts() {
		if err := component.ProvisionExternalService(); err != nil {
			return err
		}
	}

	// NOTE the code below is not tucked inside `deployment.ProvisionInstances()`
	// because I predict eventually wanting to record % progress for each instance
	// without having to expect Deployment to handle progress logic as well.

	// Concurrently provision instances
	deployment := component.ActiveDeployment()
	c := make(chan error)
	for _, instance := range deployment.Instances() {
		go func() {
			c <- instance.Provision()
		}()
	}
	for i := 0; i < component.Instances; i++ {
		if err := <-c; err != nil {
			return err
		}
	}
}
