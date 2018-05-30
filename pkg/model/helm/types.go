package helm

// Chart is representation of a helm chart.
type Chart struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ChartList stores a list of charts.
type ChartList struct {
	Items []Chart `json:"items"`
}

// Repository is representation of a helm repository.
type Repository struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Charts []Chart `json:"charts"`
}

// Repository stores a list of repositories.
type RepositoryList struct {
	Items []Repository `json:"items"`
}
