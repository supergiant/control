package kube

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type accountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type nodeProvisioner interface {
	ProvisionNodes(context.Context, []profile.NodeProfile, *model.Kube, *steps.Config) ([]string, error)
}

// Handler is a http controller for a kube entity.
type Handler struct {
	svc             Interface
	accountService  accountGetter
	nodeProvisioner nodeProvisioner
	repo            storage.Interface
}

// NewHandler constructs a Handler for kubes.
func NewHandler(svc Interface, accountService accountGetter, provisioner nodeProvisioner, repo storage.Interface) *Handler {
	return &Handler{
		svc:             svc,
		accountService:  accountService,
		nodeProvisioner: provisioner,
		repo:            repo,
	}
}

// Register adds kube handlers to a router.
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/kubes", h.createKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes", h.listKubes).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kname}", h.getKube).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kname}", h.deleteKube).Methods(http.MethodDelete)

	r.HandleFunc("/kubes/{kname}/resources", h.listResources).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kname}/resources/{resource}", h.getResource).Methods(http.MethodGet)

	r.HandleFunc("/kubes/{kname}/certs/{cname}", h.getCerts).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kname}/tasks", h.getTasks).Methods(http.MethodGet)
	// TODO(stgleb): Add get method for getting kube nodes
	r.HandleFunc("/kubes/{kname}/nodes", h.addNode).Methods(http.MethodPost)
	r.HandleFunc("/kubes/{kname}/nodes/{nodename}", h.deleteNode).Methods(http.MethodDelete)
	r.HandleFunc("/kubes/{kname}/nodes/{nodename}", h.deleteNode).Methods(http.MethodDelete)
}

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["kname"]

	if !ok {
		http.Error(w, "need name of a cluster", http.StatusBadRequest)
		return
	}

	data, err := h.repo.GetAll(r.Context(), workflows.Prefix)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to read tasks").Error(), http.StatusBadRequest)
		return
	}

	tasks := make([]*workflows.Task, 0)
	for _, v := range data {
		task := &workflows.Task{}
		err := json.Unmarshal(v, task)
		if err != nil {
			//TODO make whole handler send messages
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if task.Config.ClusterName == id {
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	cfg := tasks[0].Config

	for _, t := range tasks {
		t.Config = nil
	}

	res := &struct {
		Config *steps.Config     `json:"config"`
		Tasks  []*workflows.Task `json:"tasks"`
	}{
		Config: cfg,
		Tasks:  tasks,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) createKube(w http.ResponseWriter, r *http.Request) {
	k := &model.Kube{}
	err := json.NewDecoder(r.Body).Decode(k)
	if err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	ok, err := govalidator.ValidateStruct(k)
	if !ok {
		message.SendValidationFailed(w, err)
		return
	}

	if err = h.svc.Create(r.Context(), k); err != nil {
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]

	k, err := h.svc.Get(r.Context(), kname)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(k); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) listKubes(w http.ResponseWriter, r *http.Request) {
	kubes, err := h.svc.ListAll(r.Context())
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(kubes); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) deleteKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]
	if err := h.svc.Delete(r.Context(), kname); err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) listResources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]

	rawResources, err := h.svc.ListKubeResources(r.Context(), kname)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if _, err = w.Write(rawResources); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) getResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]
	rs := vars["resource"]
	ns := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("resourceName")

	rawResources, err := h.svc.GetKubeResources(r.Context(), kname, rs, ns, name)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if _, err = w.Write(rawResources); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) getCerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]
	cname := vars["cname"]

	b, err := h.svc.GetCerts(r.Context(), kname, cname)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(b); err != nil {
		message.SendUnknownError(w, err)
	}
}

// Add node to working kube
func (h *Handler) addNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kname := vars["kname"]
	k, err := h.svc.Get(r.Context(), kname)

	// TODO(stgleb): This method contains a lot of specific stuff, implement provision node
	// method for nodeProvisioner to do all things related to provisioning and saving cluster state
	if sgerrors.IsNotFound(err) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodeProfiles := make([]profile.NodeProfile, 0)
	err = json.NewDecoder(r.Body).Decode(&nodeProfiles)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.Get(r.Context(), k.AccountName)

	if sgerrors.IsNotFound(err) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kubeProfile := profile.Profile{
		Provider:        acc.Provider,
		Region:          k.Region,
		Arch:            k.Arch,
		OperatingSystem: k.OperatingSystem,
		UbuntuVersion:   k.OperatingSystemVersion,
		DockerVersion:   k.DockerVersion,
		K8SVersion:      k.K8SVersion,
		HelmVersion:     k.HelmVersion,

		NetworkType:    k.Networking.Type,
		CIDR:           k.Networking.CIDR,
		FlannelVersion: k.Networking.Version,

		NodesProfiles: []profile.NodeProfile{
			{},
		},

		RBACEnabled: k.RBACEnabled,
	}

	config := steps.NewConfig(k.Name, "", k.AccountName, kubeProfile)

	if len(k.Masters) != 0 {
		config.AddMaster(util.GetAnyNode(k.Masters))
	} else {
		http.Error(w, "no master found", http.StatusNotFound)
		return
	}

	// Get cloud account fill appropriate config structure with cloud account credentials
	err = util.FillCloudAccountCredentials(r.Context(), h.accountService, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)
	tasks, err := h.nodeProvisioner.ProvisionNodes(ctx, nodeProfiles, k, config)

	if err != nil && sgerrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to client side that request has been accepted
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "marshal json"))
	}
}

// TODO(stgleb): cover with unit tests
func (h *Handler) deleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kname"]
	nodeName := vars["nodename"]

	k, err := h.svc.Get(r.Context(), kname)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	// TODO(stgleb): figure out from cloud account which workflow to use
	t, err := workflows.NewTask(workflows.DigitalOceanDeleteNode, h.repo)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	config := &steps.Config{
		ClusterName:      k.Name,
		CloudAccountName: k.AccountName,
		Node: node.Node{
			Name: nodeName,
		},
	}

	err = util.FillCloudAccountCredentials(r.Context(), h.accountService, config)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	writer, err := util.GetWriter(t.ID)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	errChan := t.Run(context.Background(), *config, writer)

	// Update cluster state when deletion completes
	go func() {
		err := <-errChan

		if err != nil {
			logrus.Errorf("delete node %s from cluster %s caused %v", nodeName, kname, err)
		}

		// Delete node from cluster object
		delete(k.Nodes, nodeName)
		// Save cluster object to etcd
		logrus.Infof("delete node %s from cluster %s", nodeName, kname)
		err = h.svc.Create(context.Background(), k)

		if err != nil {
			logrus.Errorf("update cluster %s caused %v", kname, err)
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}
