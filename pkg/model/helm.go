package model

import (
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"
)

// ChartData is a simplified representation of the helm chart.
type ChartData struct {
	Metadata *chart.Metadata `json:"metadata"`
	Repo     string          `json:"repo"`
	Values   string          `json:"values"`
}

// ChartVersions is a list of the charts metadata.
type ChartVersions struct {
	Name     string              `json:"name"`
	Repo     string              `json:"repo"`
	Versions []repo.ChartVersion `json:"versions"`
}

// RepositoryInfo holds authorization details and shortened charts info.
type RepositoryInfo struct {
	Config repo.Entry      `json:"config"`
	Charts []ChartVersions `json:"charts"`
}

// ReleaseInfo is a simplified representations of the helm release.
type ReleaseInfo struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Version      int32  `json:"version"`
	CreatedAt    string `json:"createdAt"`
	LastDeployed string `json:"lastDeployed"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chartVersion"`
	Status       string `json:"status"`
}
