package digitalocean

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type CreateInstanceStep struct {
	DropletTimeout time.Duration
	CheckPeriod    time.Duration

	getServices func(string) (DropletService, KeyService)
}

func NewCreateInstanceStep(dropletTimeout, checkPeriod time.Duration) *CreateInstanceStep {
	return &CreateInstanceStep{
		DropletTimeout: dropletTimeout,
		CheckPeriod:    checkPeriod,
		getServices: func(accessToken string) (DropletService, KeyService) {
			client := digitaloceansdk.New(accessToken).GetClient()

			return client.Droplets, client.Keys
		},
	}
}

func (s *CreateInstanceStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	dropletSvc, keySvc := s.getServices(config.DigitalOceanConfig.AccessToken)
	// Node name is created from cluster name plus part of task id plus role
	config.DigitalOceanConfig.Name = util.MakeNodeName(config.ClusterName,
		config.TaskID, config.IsMaster)

	// TODO(stgleb): Move keys creation for provisioning to provisioner to be able to get
	// this key on cluster check phase.
	fingers, err := s.createKeys(ctx, keySvc, config)

	if err != nil {
		return err
	}

	tags := []string{
		config.ClusterID,
		config.DigitalOceanConfig.Name,
		config.ClusterName,
	}

	if config.IsMaster {
		tags = append(tags, fmt.Sprintf("master-%s", config.ClusterID))
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              config.DigitalOceanConfig.Name,
		Region:            config.DigitalOceanConfig.Region,
		Size:              config.DigitalOceanConfig.Size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: config.DigitalOceanConfig.Image,
		},
		Tags: tags,
	}

	role := model.RoleMaster
	if !config.IsMaster {
		role = model.RoleNode
	}

	config.Node = model.Machine{
		TaskID:   config.TaskID,
		Role:     role,
		Provider: clouds.DigitalOcean,
		Size:     config.DigitalOceanConfig.Size,
		Region:   config.DigitalOceanConfig.Region,
		State:    model.MachineStateBuilding,
		Name:     config.DigitalOceanConfig.Name,
	}

	// Update node state in cluster
	config.NodeChan() <- config.Node
	droplet, _, err := dropletSvc.Create(ctx, dropletRequest)

	if err != nil {
		config.Node.State = model.MachineStateError
		config.NodeChan() <- config.Node
		return errors.Wrap(err, "dropletService has returned an error in Run job")
	}

	after := time.After(s.DropletTimeout)
	ticker := time.NewTicker(s.CheckPeriod)

	for {
		select {
		case <-ticker.C:
			droplet, _, err = dropletSvc.Get(ctx, droplet.ID)

			if err != nil {
				return err
			}
			// Wait for droplet becomes active
			if droplet.Status == "active" {
				// Get private ip ports from droplet networks

				createdAt, _ := strconv.Atoi(droplet.Created)

				config.Node.ID = fmt.Sprintf("%d", droplet.ID)
				config.Node.CreatedAt = int64(createdAt)
				config.Node.PublicIp = getPublicIpPort(droplet.Networks.V4)
				config.Node.PrivateIp = getPrivateIpPort(droplet.Networks.V4)
				config.Node.State = model.MachineStateProvisioning
				config.Node.Name = config.DigitalOceanConfig.Name

				// Update node state in cluster
				config.NodeChan() <- config.Node

				if config.IsMaster {
					config.AddMaster(&config.Node)
				} else {
					config.AddNode(&config.Node)
				}

				logrus.Infof("Node has been created %v", config.Node)

				return nil
			}
		case <-after:
			return sgerrors.ErrTimeoutExceeded
		}
	}

	return nil
}

func (s *CreateInstanceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateInstanceStep) Name() string {
	return CreateMachineStepName
}

func (s *CreateInstanceStep) Depends() []string {
	return nil
}

func (s *CreateInstanceStep) Description() string {
	return "Create instance in Digital Ocean"
}

func (s *CreateInstanceStep) createKeys(ctx context.Context, keyService KeyService, config *steps.Config) ([]godo.DropletCreateSSHKey, error) {
	var fingers []godo.DropletCreateSSHKey

	logrus.Debugf("Step %s", CreateMachineStepName)

	// Create key for provisioning
	key, err := createKey(ctx, keyService,
		util.MakeKeyName(config.DigitalOceanConfig.Name, false),
		config.Kube.SSHConfig.BootstrapPublicKey)

	if err != nil {
		return nil, errors.Wrap(err, "create provision key")
	}

	fingers = append(fingers, godo.DropletCreateSSHKey{
		Fingerprint: key.Fingerprint,
	})

	// Create user provided key
	key, _ = createKey(ctx, keyService,
		util.MakeKeyName(config.DigitalOceanConfig.Name, true),
		config.Kube.SSHConfig.PublicKey)

	// NOTE(stgleb): In case if this key is already used by user of this account
	// just compute fingerprint and pass it
	if key == nil {
		fg, _ := fingerprint(config.Kube.SSHConfig.PublicKey)
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: fg,
		})
	} else {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: key.Fingerprint,
		})
	}

	return fingers, nil
}
