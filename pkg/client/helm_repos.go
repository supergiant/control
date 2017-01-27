package client

type HelmReposInterface interface {
	CollectionInterface
}

type HelmRepos struct {
	Collection
}
