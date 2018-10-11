package kube

type ReleaseInput struct {
	Name         string
	Namespace    string
	ChartName    string
	ChartVersion string
	RepoName     string
	Values       []byte
}
