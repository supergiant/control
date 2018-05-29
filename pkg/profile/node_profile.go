package profile

type NodeProfile struct {
	Id       string `json:"id" valid:"required"`
	Size     string `json:"size" valid:"required"`
	Image    string `json:"image" valid:"required"`
	Provider string `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	// AVX256, SSE4, MMX, AES, SR-IOV etc.
	Capabilities []string `json:"capabilities" valid:"optional"`
}
