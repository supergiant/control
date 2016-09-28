package model

type ResourceMetrics struct {
	CPUUsage int64 `json:"cpu_usage" sg:"readonly"`
	CPULimit int64 `json:"cpu_limit" sg:"readonly"`
	RAMUsage int64 `json:"ram_usage" sg:"readonly"`
	RAMLimit int64 `json:"ram_limit" sg:"readonly"`
}
