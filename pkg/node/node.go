package node

import "github.com/supergiant/supergiant/pkg/provider"

type Node struct {
	Id        string        `json:"id"`
	CreatedAt int64         `json:"created_at"`
	Provider  provider.Name `json:"provider"`
	Region    string        `json:"region"`
	PublicIp  string        `json:"public_ip"`
	PrivateIp string        `json:"private_ip"`
}
