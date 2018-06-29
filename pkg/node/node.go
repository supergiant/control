package node

type Node struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	Region    string `json:"region"`
	PublicIp  string `json:"public_ip"`
	PrivateIp string `json:"private_ip"`
}
