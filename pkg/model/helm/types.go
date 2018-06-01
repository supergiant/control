package helm

// Chart is representation of a helm chart.
type Chart struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Repository is representation of a helm repository.
type Repository struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Charts []Chart `json:"charts"`
}
