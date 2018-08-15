package provisioner

import (
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/node"
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
	req := &ProvisionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	kubeProfile, err := h.kubeService.Get(r.Context(), req.ProfileId)

	if err != nil {
		http.NotFound(w, r)
	}

	w.WriteHeader(http.StatusAccepted)
	nodeChan := make(chan *node.Node)
	errChan := make(chan error)

	for _, masterProfile := range kubeProfile.MasterProfiles {
		go func() {
			n, err := buildNode(masterProfile)

			if err != nil {
				errChan <- err
			} else {
				nodeChan <- n
			}
		}()
	}

	masterNodes := make([]*node.Node, 0, len(kubeProfile.MasterProfiles))

	for range kubeProfile.MasterProfiles {
		n := <-nodeChan
		masterNodes = append(masterNodes, n)
	}
}

func buildNode(profile profile.NodeProfile) (*node.Node, error) {
	return nil, nil
}
