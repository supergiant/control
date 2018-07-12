package node

import "github.com/supergiant/supergiant/pkg/provider"

type Node struct {
	Id        string        `json:"id" valid:"required"`
	CreatedAt int64         `json:"created_at" valid:"required"`
	Provider  provider.Name `json:"provider" valid:"required"`
	Region    string        `json:"region" valid:"required"`
	PublicIp  string        `json:"public_ip"`
	PrivateIp string        `json:"private_ip"`
}
