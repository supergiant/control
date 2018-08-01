package workflows

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/supergiant/supergiant/pkg/storage"
)

type WorkflowHandler struct {
	repository storage.Interface
}

func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of workflow", http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), "workflows", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
