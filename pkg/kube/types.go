package kube

type ReleaseInput struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ChartName    string `json:"chartName" valid:"required"`
	ChartVersion string `json:"chartVersion"`
	RepoName     string `json:"repoName" valid:"required"`
	Values       string `json:"values"`
}
