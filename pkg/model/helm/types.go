package helm

type Chart struct {
	Name string
	URL  string
}

type ChartList struct {
	Items []Chart
}

type Repository struct {
	Name   string
	URL    string
	Charts []Chart
}

type RepositoryList struct {
	Items []Repository
}