package model

type Job struct {
	ID       string `json:"-"`
	Type     int    `json:"type"`
	Data     []byte `json:"data"`
	Status   string `json:"status"`
	Attempts int    `json:"attempts"`
	Error    string `json:"error"`
}
