package model

type Addresses struct {
	External []*PortAddress `json:"external"`
	Internal []*PortAddress `json:"internal"`
}

type PortAddress struct {
	Port    string `json:"port"` // TODO really this should be the name of the port, which currently is the string of the number
	Address string `json:"address"`
}

// NOTE this is not to be confused with our concept of Resources like Apps and
// Components -- this is for CPU / RAM / disk.
type ResourceMetrics struct {
	CPUUsage int64 `json:"cpu_usage" sg:"readonly"`
	CPULimit int64 `json:"cpu_limit" sg:"readonly"`
	RAMUsage int64 `json:"ram_usage" sg:"readonly"`
	RAMLimit int64 `json:"ram_limit" sg:"readonly"`
}

type Pagination struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}
