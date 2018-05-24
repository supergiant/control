package cloud_account

import (
	"encoding/json"
	"net/http"
)

type Endpoint struct {
	Service Service
}

func (e *Endpoint) CreateCloudAccount(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	err := json.NewDecoder(r.Body).Decode(account)

	//TODO Revise error handling
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = e.Service.Create(r.Context(), account); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}
