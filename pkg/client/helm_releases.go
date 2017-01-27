package client

type HelmReleasesInterface interface {
	CollectionInterface
}

type HelmReleases struct {
	Collection
}
