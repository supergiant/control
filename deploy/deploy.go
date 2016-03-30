package deploy

import "github.com/supergiant/supergiant/client"

func Deploy(appName string, componentName string, currentReleaseID string, targetReleaseID string) error {

	sg := client.New("http://localhost:8080")

	app, err := sg.Apps().Get(appName)
	if err != nil {
		return err
	}

	component, err := app.Components().Get(componentName)
	if err != nil {
		return err
	}

	currentRelease, err := component.CurrentRelease()
	if err != nil {
		return err
	}

	targetRelease, err := component.TargetRelease()
	if err != nil {
		return err
	}

	targetList, err := targetRelease.Instances().List()
	if err != nil {
		return err
	}
	targetInstances := targetList.Items

	if currentRelease == nil { // first release
		for _, instance := range targetInstances {
			if err = instance.Start(); err != nil {
				return err
			}
		}
		for _, instance := range targetInstances {
			if err = instance.WaitForStarted(); err != nil {
				return err
			}
		}
		return nil
	}

	currentList, err := currentRelease.Instances().List()
	if err != nil {
		return err
	}
	currentInstances := currentList.Items

	// remove instances
	if currentRelease.InstanceCount > targetRelease.InstanceCount {
		instancesRemoving := currentRelease.InstanceCount - targetRelease.InstanceCount
		for _, instance := range currentInstances[len(currentInstances)-instancesRemoving:] {
			if err := instance.Stop(); err != nil {
				return err
			}
		}
	}

	// update instances
	// NOTE instance count should be the same at this point between releases
	for i := 0; i < currentRelease.InstanceCount; i++ {
		currentInstance := currentInstances[i]
		targetInstance := targetInstances[i]

		currentInstance.Stop()
		currentInstance.WaitForStopped()

		targetInstance.Start()
		targetInstance.WaitForStarted()
	}

	// add new instances
	if currentRelease.InstanceCount < targetRelease.InstanceCount {
		instancesAdding := targetRelease.InstanceCount - currentRelease.InstanceCount
		newInstances := targetInstances[len(targetInstances)-instancesAdding:]
		for _, instance := range newInstances {
			if err := instance.Start(); err != nil {
				return err
			}
		}
		for _, instance := range newInstances {
			if err := instance.WaitForStarted(); err != nil {
				return err
			}
		}
	}

	return nil
}
