package model

import (
	"fmt"

	"github.com/supergiant/control/pkg/clouds"
)

type MachineState string

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	MachineStatePlanned      MachineState = "planned"
	MachineStateBuilding     MachineState = "building"
	MachineStateProvisioning MachineState = "provisioning"
	MachineStateError        MachineState = "error"
	MachineStateActive       MachineState = "active"
	MachineStateDeleting     MachineState = "deleting"

	RoleMaster Role = "master"
	RoleNode   Role = "node"
)

type Machine struct {
	ID               string       `json:"id" valid:"required"`
	TaskID           string       `json:"taskId"`
	Role             Role         `json:"role"`
	CreatedAt        int64        `json:"createdAt" valid:"required"`
	Provider         clouds.Name  `json:"provider" valid:"required"`
	Region           string       `json:"region" valid:"required"`
	AvailabilityZone string       `json:"az" valid:"-"`
	Size             string       `json:"size"`
	PublicIp         string       `json:"publicIp"`
	PrivateIp        string       `json:"privateIp"`
	State            MachineState `json:"state"`
	Name             string       `json:"name"`
}

func (m Machine) String() string {
	return fmt.Sprintf("<ID: %s, Name: %s Active: %v, Size: %s, CreatedAt: %d, Provider: %s, Region; %s, AvailabilityZone: %s, PublicIp: %s, PrivateIp: %s>",
		m.ID, m.Name, m.State, m.Size, m.CreatedAt, m.Provider, m.Region,
		m.AvailabilityZone, m.PublicIp, m.PrivateIp)
}

func ToRole(isMaster bool) Role {
	if isMaster {
		return RoleMaster
	}
	return RoleNode
}
