package node

import (
	"fmt"

	"github.com/supergiant/supergiant/pkg/clouds"
)

// TODO(stgleb): Accommodate terminology and rename Node to Machine
type Node struct {
	Id        string      `json:"id" valid:"required"`
	Role      string      `json:"role"`
	CreatedAt int64       `json:"created_at" valid:"required"`
	Provider  clouds.Name `json:"provider" valid:"required"`
	Region    string      `json:"region" valid:"required"`
	Size      string      `json:"size"`
	PublicIp  string      `json:"public_ip"`
	PrivateIp string      `json:"private_ip"`
	Active    bool        `json:"active"`
	Name      string      `json:"name"`
}

func (n Node) String() string {
	return fmt.Sprintf("<ID: %s, Active: %v, Size: %s, CreatedAt: %d, Provider: %s, Region; %s, PublicIp: %s, PrivateIp: %s>",
		n.Id, n.Active, n.Size, n.CreatedAt, n.Provider, n.Region, n.PublicIp, n.PrivateIp)
}
