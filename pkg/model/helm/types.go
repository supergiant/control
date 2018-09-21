package helm

import (
	"time"

	"k8s.io/helm/pkg/repo"
)

// Chart is representation of a helm chart.
type Chart struct {
	Name          string         `json:"name"`
	Repo          string         `json:"repo"`
	Description   string         `json:"description"`
	Home          string         `json:"home"`
	Keywords      []string       `json:"keywords"`
	Maintainers   []Maintainer   `json:"maintainers"`
	Sources       []string       `json:"sources"`
	Icon          string         `json:"icon"`
	ChartVersions []ChartVersion `json:"chartVersions"`
}

type ChartVersion struct {
	Version    string
	AppVersion string
	Created    time.Time
	Digest     string
	URLs       []string
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
