package model

import (
	"time"

	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"
)

// ChartData is a simplified representation of the helm chart.
type ChartData struct {
	Metadata *chart.Metadata `json:"metadata"`
	Readme   string          `json:"readme"`
	Values   string          `json:"values"`
}

// ChartInfo is a list of the charts metadata.
type ChartInfo struct {
	Name        string         `json:"name"`
	Repo        string         `json:"repo"`
	Description string         `json:"description"`
	Versions    []ChartVersion `json:"versions"`
}

type ChartVersion struct {
	Version    string    `json:"version"`
	AppVersion string    `json:"appVersion"`
	Created    time.Time `json:"created"`
	Digest     string    `json:"digest"`
	URLs       []string  `json:"urls"`
}

// RepositoryInfo holds authorization details and shortened charts info.
type RepositoryInfo struct {
	Config repo.Entry  `json:"config"`
	Charts []ChartInfo `json:"charts"`
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
