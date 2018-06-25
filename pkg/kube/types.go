package kube

// Kube represents a kubernetes cluster.
type Kube struct {
	Name        string `json:"name" valid:"required"`
	Version     string `json:"version"`
	RBACEnabled bool   `json:"rbacEnabled"`
	AccountName string `json:"accountName"`
	APIAddr     string `json:"apiAddr"`
	Auth        Auth   `json:"auth"`
	SSHUser     string `json:"sshUser"`
	SSHKey      []byte `json:"sshKey"`
}

// Auth holds all possible auth parameters.
type Auth struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	CA       string `json:"ca"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
