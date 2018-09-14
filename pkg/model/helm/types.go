package helm

import "k8s.io/helm/pkg/repo"

// Chart is representation of a helm chart.
type Chart struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Repository is representation of a helm repository.
type Repository struct {
	Config repo.Entry     `json:"config"`
	Index  repo.IndexFile `json:"index"`
}
