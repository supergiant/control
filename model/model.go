package model

// TODO
type Resource interface {
	EtcdKey(id string) string
	// NewModel() Model
}

type Model interface {
}

// type ModelList struct {
//   Items []Model
// }
