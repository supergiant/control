package model

import (
	"time"

	"k8s.io/helm/pkg/repo"
)

// Chart is representation of a helm chart.
type Chart struct {
	Name        string         `json:"name"`
	Repo        string         `json:"repo"`
	Description string         `json:"description"`
	Home        string         `json:"home"`
	Keywords    []string       `json:"keywords"`
	Maintainers []Maintainer   `json:"maintainers"`
	Sources     []string       `json:"sources"`
	Icon        string         `json:"icon"`
	Versions    []ChartVersion `json:"versions"`
}

type ChartVersion struct {
	Version    string    `json:"version"`
	AppVersion string    `json:"appVersion"`
	Created    time.Time `json:"created"`
	Digest     string    `json:"digest"`
	URLs       []string  `json:"urls"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository is representation of a helm repository.
type Repository struct {
	Config repo.Entry `json:"config"`
	Charts []Chart    `json:"charts"`
}

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
