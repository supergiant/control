package provisioner

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/profile"
)

type ProvisionHandler struct {
	kubeService profile.KubeProfileService
	nodeService profile.NodeProfileService
}

type ProvisionRequest struct {
	ProfileId string `json:"profileId"`
}

func New(kubeService profile.KubeProfileService, nodeService profile.NodeProfileService) *ProvisionHandler {
	return &ProvisionHandler{
		kubeService: kubeService,
		nodeService: nodeService,
	}
}

func (h *ProvisionHandler) Provision(w http.ResponseWriter, r *http.Request) {

}
