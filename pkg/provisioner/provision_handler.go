package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
)

type ProvisionHandler struct {
	kubeService profile.KubeProfileService
	nodeService profile.NodeProfileService
	repository  storage.Interface
}

type ProvisionRequest struct {
	ProfileId string `json:"profileId"`
}

type ProvisionResponse struct {
	TaskIds []string `json:"taskIds"`
}

func NewHandler(kubeService profile.KubeProfileService, nodeService profile.NodeProfileService) *ProvisionHandler {
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

	// Respond to client side that request has been accepted
	w.WriteHeader(http.StatusAccepted)
	provisioner := NewProvisioner(h.repository)
	taskIds := provisioner.Prepare(len(kubeProfile.MasterProfiles), len(kubeProfile.NodesProfiles))

	resp := ProvisionResponse{
		TaskIds: taskIds,
	}

	err = json.NewEncoder(w).Encode(&resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO(stgleb): Cover infrastructure creation in task
	nodeChan := make(chan node.Node)
	errChan := make(chan error)

	profiles := append(append(make([]profile.NodeProfile, 0),
		kubeProfile.MasterProfiles...),
		kubeProfile.NodesProfiles...)
	nodes := make([]node.Node, 0, len(kubeProfile.MasterProfiles)+len(kubeProfile.NodesProfiles))

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

	provisioner.Provision(context.Background(), nodes)
}

func buildNode(profile profile.NodeProfile) (node.Node, error) {
	return node.Node{}, nil
}
