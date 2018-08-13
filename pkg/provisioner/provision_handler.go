package provisioner

import (
	"net/http"
	"github.com/supergiant/supergiant/pkg/profile"
)

type ProvisionHandler struct{

}

type ProvisionRequest struct{
	Profile profile.KubeProfile	`json:"profile"`
}

func (h *ProvisionHandler) Provision(w http.ResponseWriter,r *http.Request) {

}