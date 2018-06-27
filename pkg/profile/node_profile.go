package profile

import "github.com/supergiant/supergiant/pkg/provider"

type NodeProfile struct {
	ID       string        `json:"id" valid:"required"`
	Size     *NodeSize     `json:"size" valid:"required"`
	Image    string        `json:"image" valid:"required"`
	Provider provider.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	// AVX256, SSE4, MMX, AES, SR-IOV etc.
	Capabilities []string `json:"capabilities" valid:"optional"`
	Labels       []string `json:"labels"`
}

type NodeSize struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`

	// Name e.g. t2.micro
	Name string `json:"name"`
}
