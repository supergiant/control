package node

import (
	"fmt"

	"github.com/supergiant/supergiant/pkg/clouds"
)

type NodeState string

type Role string

const (
	StatePlanned      NodeState = "planned"
	StateBuilding     NodeState = "building"
	StateProvisioning NodeState = "provisioning"
	StateError        NodeState = "error"
	StateActive       NodeState = "active"

	RoleMaster Role = "master"
	RoleNode   Role = "node"
)

// TODO(stgleb): Accommodate terminology and rename Node to Machine
type Node struct {
	ID        string      `json:"id" valid:"required"`
	TaskID    string      `json:"taskId"`
	Role      Role        `json:"role"`
	CreatedAt int64       `json:"createdAt" valid:"required"`
	Provider  clouds.Name `json:"provider" valid:"required"`
	Region    string      `json:"region" valid:"required"`
	Size      string      `json:"size"`
	PublicIp  string      `json:"publicIp"`
	PrivateIp string      `json:"privateIp"`
	State     NodeState   `json:"state"`
	Name      string      `json:"name"`
}

func (n Node) String() string {
	return fmt.Sprintf("<ID: %s, Name: %s Active: %v, Size: %s, CreatedAt: %d, Provider: %s, Region; %s, PublicIp: %s, PrivateIp: %s>",
		n.ID, n.Name, n.State, n.Size, n.CreatedAt, n.Provider, n.Region, n.PublicIp, n.PrivateIp)
}
