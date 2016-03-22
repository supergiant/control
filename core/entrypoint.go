package core

type Entrypoint struct {
	Domain  string `json:"domain"`  // e.g. blog.qbox.io
	Address string `json:"address"` // the ELB address

	// NOTE we actually don't need this -- we can always attach the policy, and enable per port
	// IPWhitelistEnabled bool   `json:"ip_whitelist_enabled"`
}
