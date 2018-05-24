package profile

type NodeProfile struct {
	Id       string `json:"id"`
	Size     string `json:"size"`
	Image    string `json:"image"`
	Provider string `json:"provider"`
	// AVX256, SSE4, MMX, AES, SR-IOV etc.
	Capabilities []string `json:"capabilities"`
}
