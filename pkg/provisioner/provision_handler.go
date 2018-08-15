package provisioner

import (
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"context"
)

type ProvisionHandler struct {
	kubeService profile.KubeProfileService
	nodeService profile.NodeProfileService
	provisioner Provisioner
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
	nodeChan := make(chan node.Node)
	errChan := make(chan error)

	profiles := append(append(make([]profile.NodeProfile, 0),
		kubeProfile.MasterProfiles...),
		kubeProfile.NodesProfiles...)
	nodes := make([]node.Node, 0, len(kubeProfile.MasterProfiles) + len(kubeProfile.NodesProfiles))

	// Build master nodes
	for _, p := range profiles {
		go func() {
			n, err := buildNode(p)

			if err != nil {
				errChan <- err
			} else {
				nodeChan <- n
			}
		}()
	}

	for i := 0; i < len(profiles); i++ {
		n := <-nodeChan
		nodes = append(nodes, n)
	}

	h.provisioner.Provision(context.Background(), nodes)
}

func buildNode(profile profile.NodeProfile) (node.Node, error) {
	return node.Node{}, nil
}
