package node

import (
	"github.com/supergiant/supergiant/pkg/clouds"
	"fmt"
)

type Node struct {
	Id        string      `json:"id" valid:"required"`
	CreatedAt int64       `json:"created_at" valid:"required"`
	Provider  clouds.Name `json:"provider" valid:"required"`
	Region    string      `json:"region" valid:"required"`
	PublicIp  string      `json:"public_ip"`
	PrivateIp string      `json:"private_ip"`
}


func (n Node) String() string {
	return fmt.Sprintf("<Id: %s, CreatedAt: %d, Provider: %s, Region; %s, PublicIp: %s, PrivateIp: %s>",
		n.Id, n.CreatedAt, n.Provider, n.Region, n.PublicIp, n.PrivateIp)
}