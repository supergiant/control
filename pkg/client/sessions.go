package client

type SessionsInterface interface {
	CollectionInterface
}

type Sessions struct {
	Collection
}

// NOTE update will give 404
