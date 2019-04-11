package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/proxy"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	clusterService = "kubernetes.io/cluster-service"
	nodeLabelRole  = "kubernetes.io/role"
)

type accountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type profileGetter interface {
	Get(context.Context, string) (*profile.Profile, error)
}

type nodeProvisioner interface {
	ProvisionNodes(context.Context, []profile.NodeProfile, *model.Kube,
		*steps.Config) ([]string, error)
	// Method that cancels newly added nodes to working cluster
	Cancel(string) error
}

type kubeProvisioner interface {
	RestartClusterProvisioning(ctx context.Context,
		clusterProfile *profile.Profile,
		config *steps.Config,
		taskIdMap map[string][]string) error
}

type ServiceInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Namespace string `json:"namespace"`
	ProxyPort string `json:"proxyPort"`
}

type MetricResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// Handler is a http controller for a kube entity.
type Handler struct {
	svc             Interface
	accountService  accountGetter
	nodeProvisioner nodeProvisioner
	kubeProvisioner kubeProvisioner
	profileSvc      profileGetter

	repo    storage.Interface
	proxies proxy.Container

	getWriter       func(string) (io.WriteCloser, error)
	getMetrics      func(string, *model.Kube) (*MetricResponse, error)
	listK8sServices func(*model.Kube, string) (*corev1.ServiceList, error)
}

// NewHandler constructs a Handler for kubes.
func NewHandler(
	svc Interface,
	accountService accountGetter,
	profileSvc profileGetter,
	provisioner nodeProvisioner,
	kubeProvisioner kubeProvisioner,
	repo storage.Interface,
	proxies proxy.Container,
) *Handler {
	return &Handler{
		svc:             svc,
		accountService:  accountService,
		nodeProvisioner: provisioner,
		kubeProvisioner: kubeProvisioner,
		profileSvc:      profileSvc,
		repo:            repo,
		getWriter:       util.GetWriter,
		getMetrics: func(metricURI string, k *model.Kube) (*MetricResponse, error) {
			cfg, err := NewConfigFor(k)
			if err != nil {
				return nil, errors.Wrap(err, "build kubernetes rest config")
			}
			kclient, err := rest.UnversionedRESTClientFor(cfg)
			if err != nil {
				return nil, errors.Wrap(err, "build kubernetes client")
			}

			raw, err := kclient.Get().RequestURI(metricURI).Do().Raw()
			if err != nil {
				return nil, errors.Wrap(err, "retrieve metrics")
			}

			metricResponse := &MetricResponse{}
			err = json.Unmarshal(raw, metricResponse)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal")
			}

			return metricResponse, nil
		},
		listK8sServices: func(k *model.Kube, selector string) (*corev1.ServiceList, error) {
			cfg, err := NewConfigFor(k)
			if err != nil {
				return nil, errors.Wrap(err, "build kubernetes rest config")
			}
			c, err := clientcorev1.NewForConfig(cfg)
			if err != nil {
				return nil, errors.Wrapf(err, "build kubernetes client")
			}
			return c.Services(metav1.NamespaceAll).List(metav1.ListOptions{
				LabelSelector: selector,
			})
		},
		proxies: proxies,
	}
}

// Register adds kube handlers to a router.
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/kubes", h.createKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes", h.listKubes).Methods(http.MethodGet)
	r.HandleFunc("/kubes/import", h.importKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes/{kubeID}", h.getKube).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}", h.deleteKube).Methods(http.MethodDelete)

	r.HandleFunc("/kubes/{kubeID}/users/{uname}/kubeconfig", h.getKubeconfig).Methods(http.MethodGet)

	r.HandleFunc("/kubes/{kubeID}/resources", h.listResources).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/resources/{resource}", h.getResource).Methods(http.MethodGet)

	r.HandleFunc("/kubes/{kubeID}/releases", h.installRelease).Methods(http.MethodPost)
	r.HandleFunc("/kubes/{kubeID}/releases", h.listReleases).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/releases/{releaseName}", h.getRelease).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/releases/{releaseName}", h.deleteReleases).Methods(http.MethodDelete)

	r.HandleFunc("/kubes/{kubeID}/certs/{cname}", h.getCerts).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/tasks", h.getTasks).Methods(http.MethodGet)

	// DEPRECATED: has been moved to /kubes/{kubeID}/machines
	r.HandleFunc("/kubes/{kubeID}/nodes", h.addMachine).Methods(http.MethodPost)
	// DEPRECATED: has been moved to /kubes/{kubeID}/machines
	r.HandleFunc("/kubes/{kubeID}/nodes/{nodename}", h.deleteMachine).Methods(http.MethodDelete)

	r.HandleFunc("/kubes/{kubeID}/nodes", h.listNodes).Methods(http.MethodGet)

	r.HandleFunc("/kubes/{kubeID}/machines", h.addMachine).Methods(http.MethodPost)
	r.HandleFunc("/kubes/{kubeID}/machines/{nodename}", h.deleteMachine).Methods(http.MethodDelete)

	r.HandleFunc("/kubes/{kubeID}/nodes/metrics", h.getNodesMetrics).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/metrics", h.getClusterMetrics).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/services", h.getServices).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/restart", h.restartKubeProvisioning).Methods(http.MethodPost)
}

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["kubeID"]

	if !ok {
		http.Error(w, "need name of a cluster", http.StatusBadRequest)
		return
	}

	tasks, err := h.getKubeTasks(r.Context(), id)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, id, err)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	if len(tasks) == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	type taskDTO struct {
		ID           string                 `json:"id"`
		Type         string                 `json:"type"`
		Status       statuses.Status        `json:"status"`
		StepStatuses []workflows.StepStatus `json:"stepsStatuses"`
	}

	resp := make([]taskDTO, 0, len(tasks))

	for _, task := range tasks {
		resp = append(resp, taskDTO{
			ID:           task.ID,
			Type:         task.Type,
			Status:       task.Status,
			StepStatuses: task.StepStatuses,
		})
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) createKube(w http.ResponseWriter, r *http.Request) {
	newKube := &model.Kube{}
	err := json.NewDecoder(r.Body).Decode(newKube)
	if err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	ok, err := govalidator.ValidateStruct(newKube)
	if !ok {
		message.SendValidationFailed(w, err)
		return
	}

	existingKube, err := h.svc.Get(r.Context(), newKube.ID)
	if existingKube != nil {
		message.SendAlreadyExists(w, existingKube.ID, sgerrors.ErrAlreadyExists)
		return
	}

	if err != nil && !sgerrors.IsNotFound(err) {
		message.SendUnknownError(w, err)
		return
	}

	if err = h.svc.Create(r.Context(), newKube); err != nil {
		message.SendUnknownError(w, err)
		return
	}

	// TODO(stgleb): Reply with kube ID
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
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
	kubeID := vars["kubeID"]
	logrus.Debugf("Delete kube %s", kubeID)

	if err := h.nodeProvisioner.Cancel(kubeID); err != nil {
		logrus.Debugf("cancel kube tasks error %v", err)
	}

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	acc, err := h.accountService.Get(r.Context(), k.AccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	config := &steps.Config{
		Provider:         k.Provider,
		ClusterID:        k.ID,
		ClusterName:      k.Name,
		CloudAccountName: k.AccountName,
		Masters:          steps.NewMap(k.Masters),
		Nodes:            steps.NewMap(k.Nodes),
	}

	t, err := workflows.NewTask(config, workflows.DeleteCluster, h.repo)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	// Load things specific to cloud provider
	err = util.LoadCloudSpecificDataFromKube(k, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	err = util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	fileName := util.MakeFileName(t.ID)
	writer, err := h.getWriter(fileName)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	errChan := t.Run(context.Background(), *config, writer)

	go func(t *workflows.Task) {
		// Update kube with deleting state
		k.State = model.StateDeleting
		err = h.svc.Create(context.Background(), k)

		if err != nil {
			logrus.Errorf("update cluster %s caused %v", kubeID, err)
		}

		err = <-errChan
		if err != nil {
			return
		}

		// Clean up tasks in storage
		err = h.deleteClusterTasks(context.Background(), kubeID)

		if err != nil {
			logrus.Errorf("error while deleting tasks %s", err)
		}

		// Finally delete cluster record from etcd
		if err := h.svc.Delete(context.Background(), kubeID); err != nil {
			logrus.Errorf("delete kube %s caused %v", kubeID, err)
			return
		}
	}(t)

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getKubeconfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["kubeID"]
	user := vars["uname"]

	data, err := h.svc.KubeConfigFor(r.Context(), kname, user)
	if err != nil {
		logrus.Errorf("kubes: %s cluster: get kubeconfig: %s", kname, err)
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, user, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if _, err = w.Write(data); err != nil {
		logrus.Errorf("kubes: %s cluster: get kubeconfig: write response: %s", kname, err)
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) listResources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]
	rawResources, err := h.svc.ListKubeResources(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
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

	kubeID := vars["kubeID"]
	rs := vars["resource"]
	ns := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	rawResources, err := h.svc.GetKubeResources(r.Context(), kubeID, rs, ns, name)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
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

	kubeID := vars["kubeID"]
	cname := vars["cname"]

	b, err := h.svc.GetCerts(r.Context(), kubeID, cname)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(b); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) listNodes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeID := vars["kubeID"]
	role := r.URL.Query().Get("role")

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	nodes, err := h.svc.ListNodes(r.Context(), k, role)
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(nodes); err != nil {
		message.SendUnknownError(w, err)
	}
}

// Add node to working kube
func (h *Handler) addMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeID := vars["kubeID"]
	k, err := h.svc.Get(r.Context(), kubeID)

	if sgerrors.IsNotFound(err) {
		http.NotFound(w, r)
		return
	}

	logrus.Debugf("Get cloud profile %s", k.ProfileID)
	kubeProfile, err := h.profileSvc.Get(r.Context(), k.ProfileID)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, k.ProfileID, err)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	config, err := steps.NewConfigFromKube(kubeProfile, k)
	if err != nil {
		logrus.Errorf("New config %v", err.Error())
		message.SendUnknownError(w, err)
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

	if len(k.Masters) != 0 {
		config.AddMaster(util.GetRandomNode(k.Masters))
	} else {
		http.Error(w, "no master found", http.StatusNotFound)
		return
	}

	// Get cloud account fill appropriate config structure
	// with cloud account credentials
	err = util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*20)
	tasks, err := h.nodeProvisioner.ProvisionNodes(ctx, nodeProfiles,
		k, config)

	if err != nil && sgerrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add tasks ids to kube object
	k.Tasks[workflows.NodeTask] = append(k.Tasks[workflows.NodeTask], tasks...)

	if err := h.svc.Create(ctx, k); err != nil {
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
func (h *Handler) deleteMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]
	nodeName := vars["nodename"]

	logrus.Debugf("Delete node %s from kube %s",
		nodeName, kubeID)
	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	// TODO(stgleb): check whether we will have quorum of master nodes if node is deleted.
	if _, ok := k.Masters[nodeName]; ok {
		http.Error(w, "delete master node not allowed", http.StatusMethodNotAllowed)
		return
	}

	var n *model.Machine

	if n = k.Nodes[nodeName]; n == nil {
		http.NotFound(w, r)
		return
	}

	acc, err := h.accountService.Get(r.Context(), k.AccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	config := &steps.Config{
		Kube:     *k,
		Provider: k.Provider,
		DrainConfig: steps.DrainConfig{
			PrivateIP: n.PrivateIp,
		},
		ClusterID:        k.ID,
		ClusterName:      k.Name,
		CloudAccountName: k.AccountName,
		Node:             *n,
		Masters:          steps.NewMap(k.Masters),
	}

	t, err := workflows.NewTask(config, workflows.DeleteNode, h.repo)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	err = util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	err = util.LoadCloudSpecificDataFromKube(k, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	writer, err := h.getWriter(util.MakeFileName(t.ID))

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	// Update cluster state when deletion completes
	go func() {
		// Set node to deleting state
		nodeToDelete, ok := k.Nodes[nodeName]

		if !ok {
			logrus.Errorf("Node %s not found", nodeName)
			return
		}
		nodeToDelete.State = model.MachineStateDeleting
		k.Nodes[nodeName] = nodeToDelete
		err := h.svc.Create(context.Background(), k)

		if err != nil {
			logrus.Errorf("update cluster %s caused %v", kubeID, err)
		}

		err = <-t.Run(context.Background(), *config, writer)

		if err != nil {
			logrus.Errorf("delete node %s from cluster %s caused %v", nodeName, kubeID, err)
		}

		// Delete node from cluster object
		delete(k.Nodes, nodeName)
		// Save cluster object to etcd
		logrus.Infof("delete node %s from cluster %s", nodeName, kubeID)
		err = h.svc.Create(context.Background(), k)

		if err != nil {
			logrus.Errorf("update cluster %s caused %v", kubeID, err)
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}

// TODO(stgleb): Create separte task service to manage task object lifecycle
func (h *Handler) getKubeTasks(ctx context.Context, kubeID string) ([]*workflows.Task, error) {
	k, err := h.svc.Get(ctx, kubeID)

	if err != nil {
		return nil, err
	}

	tasks := make([]*workflows.Task, 0, len(k.Tasks))

	for _, taskSet := range k.Tasks {
		for _, taskID := range taskSet {
			t, err := h.repo.Get(ctx, workflows.Prefix, taskID)

			// If one of tasks not found we dont care, because
			// they may npt be created yet
			if err != nil {
				logrus.Debugf("task %s not found", taskID)
				continue
			}

			task := &workflows.Task{}
			err = json.Unmarshal(t, task)

			if err != nil {
				return nil, errors.Wrapf(err,
					"get task %s", taskID)
			}

			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (h *Handler) deleteClusterTasks(ctx context.Context, kubeID string) error {
	tasks, err := h.getKubeTasks(ctx, kubeID)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete cluster %s tasks", kubeID))
	}

	for _, task := range tasks {
		if err := h.repo.Delete(ctx, workflows.Prefix, task.ID); err != nil {
			logrus.Warnf("delete task %s: %v", task.ID, err)
			return err
		}
	}

	return nil
}

func (h *Handler) installRelease(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	inp := &ReleaseInput{}
	err := json.NewDecoder(r.Body).Decode(inp)
	if err != nil {
		logrus.Errorf("helm: install release: decode: %s", err)
		message.SendInvalidJSON(w, err)
		return
	}
	ok, err := govalidator.ValidateStruct(inp)
	if !ok {
		logrus.Errorf("helm: install release: validation: %s", err)
		message.SendValidationFailed(w, err)
		return
	}

	kubeID := vars["kubeID"]
	rls, err := h.svc.InstallRelease(r.Context(), kubeID, inp)
	if err != nil {
		logrus.Errorf("helm: install release: %s cluster: %s (%+v)", kubeID, err, inp)
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(rls); err != nil {
		logrus.Errorf("helm: install release: %s cluster: %s/%s: write response: %s",
			kubeID, inp.RepoName, inp.ChartName, err)
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) getRelease(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]
	rlsName := vars["releaseName"]

	rls, err := h.svc.ReleaseDetails(r.Context(), kubeID, rlsName)
	if err != nil {
		logrus.Errorf("helm: get %s release: %s cluster: %s", rlsName, kubeID, err)
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(rls); err != nil {
		logrus.Errorf("helm: get %s release: %s cluster: write response: %s", rlsName, kubeID, err)
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) listReleases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]
	// TODO: use a struct for input parameters
	rlsList, err := h.svc.ListReleases(r.Context(), kubeID, "", "", 0)
	if err != nil {
		logrus.Errorf("helm: list releases: %s cluster: %s", kubeID, err)
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(rlsList); err != nil {
		logrus.Errorf("helm: list releases: %s cluster: write response: %s", kubeID, err)
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) deleteReleases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kubeID := vars["kubeID"]
	rlsName := vars["releaseName"]
	purge, _ := strconv.ParseBool(r.URL.Query().Get("purge"))

	rls, err := h.svc.DeleteRelease(r.Context(), kubeID, rlsName, purge)
	if err != nil {
		logrus.Errorf("helm: delete release: %s cluster: release %s: %s", kubeID, rlsName, err)
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(rls); err != nil {
		logrus.Errorf("helm: delete release: %s cluster: write response: %s", kubeID, err)
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) getClusterMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		metricsRelUrls = map[string]string{
			"cpu":    "api/v1/query?query=:node_cpu_utilisation:avg1m",
			"memory": "api/v1/query?query=:node_memory_utilisation:",
		}
		masterNode *model.Machine
		response   = map[string]interface{}{}
		baseUrl    = "api/v1/namespaces/kube-system/services/prometheus-operated:9090/proxy"
	)

	vars := mux.Vars(r)
	kubeID := vars["kubeID"]

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	for key := range k.Masters {
		if k.Masters[key] != nil {
			masterNode = k.Masters[key]
		}
	}

	if masterNode == nil {
		return
	}

	for metricType, relUrl := range metricsRelUrls {
		url := fmt.Sprintf("https://%s/%s/%s", masterNode.PublicIp, baseUrl, relUrl)
		metricResponse, err := h.getMetrics(url, k)

		if err != nil {
			message.SendUnknownError(w, err)
			return
		}

		if len(metricResponse.Data.Result) > 0 && len(metricResponse.Data.Result[0].Value) > 1 {
			response[metricType] = metricResponse.Data.Result[0].Value[1]
		}
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) getNodesMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		metricsRelUrls = map[string]string{
			"cpu":    "api/v1/query?query=node:node_cpu_utilisation:avg1m",
			"memory": "api/v1/query?query=node:node_memory_utilisation:",
		}
		masterNode *model.Machine
		response   = map[string]map[string]interface{}{}
		baseUrl    = "api/v1/namespaces/kube-system/services/prometheus-operated:9090/proxy"
	)

	vars := mux.Vars(r)
	kubeID := vars["kubeID"]

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	for key := range k.Masters {
		if k.Masters[key] != nil {
			masterNode = k.Masters[key]
		}
	}

	if masterNode == nil {
		return
	}

	for metricType, relUrl := range metricsRelUrls {
		url := fmt.Sprintf("https://%s/%s/%s", masterNode.PublicIp, baseUrl, relUrl)
		metricResponse, err := h.getMetrics(url, k)

		if err != nil {
			message.SendUnknownError(w, err)
			return
		}

		for _, result := range metricResponse.Data.Result {
			// Get node name of the metric
			nodeName, ok := result.Metric["node"]

			if !ok {
				continue
			}
			// If dict for this node is empty - fill it with empty map
			if response[nodeName] == nil {
				response[nodeName] = map[string]interface{}{}
			}

			response[nodeName][metricType] = result.Value[1]
		}
	}

	if k.Provider == clouds.AWS {
		processAWSMetrics(k, response)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) getServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeID := vars["kubeID"]

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	// TODO: use sg specific label
	selector := fmt.Sprintf("%s=%s", clusterService, "true")
	svcList, err := h.listK8sServices(k, selector)
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	// TODO(stgleb): Figure out which ports are worth to be proxy. These are name aliases for ports!
	webPorts := map[string]struct{}{
		"web":     {},
		"http":    {},
		"https":   {},
		"service": {},
	}

	var serviceInfos = make([]*ServiceInfo, 0)
	var targetServices = make([]proxy.Target, 0)

	cfg, err := NewConfigFor(k)
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}
	for _, service := range svcList.Items {
		for _, port := range service.Spec.Ports {
			if _, ok := webPorts[port.Name]; !ok && port.Protocol != "TCP" {
				continue
			}

			var serviceInfo = &ServiceInfo{
				ID:        string(service.UID),
				Name:      service.Name,
				Namespace: service.Namespace,
				Type:      string(service.Spec.Type),
			}
			serviceInfos = append(serviceInfos, serviceInfo)

			targetServices = append(targetServices, proxy.Target{
				ProxyID: kubeID + string(service.ObjectMeta.UID),
				// TargetPort is ignored here
				TargetURL:  fmt.Sprintf("%s%s:%d/proxy", cfg.Host, service.ObjectMeta.SelfLink, port.Port),
				KubeConfig: cfg,
			})
		}
	}

	// TODO implement proxy removing logic under separate ticket
	err = h.proxies.RegisterProxies(targetServices)
	if err != nil {
		logrus.Error(err)
		message.SendUnknownError(w, err)
		return
	}

	proxies := h.proxies.GetProxies(kubeID)
	for _, service := range serviceInfos {
		if proxies[kubeID+service.ID] == nil {
			continue
		}
		service.ProxyPort = proxies[kubeID+service.ID].Port()
	}

	err = json.NewEncoder(w).Encode(serviceInfos)
	if err != nil {
		message.SendUnknownError(w, err)
	}
}

func contains(name, value string, labels map[string]string) bool {
	v, exists := labels[name]
	if exists && v == value {
		return true
	}

	return false
}

func (h *Handler) restartKubeProvisioning(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeID := vars["kubeID"]

	logrus.Debugf("Get kube %s", kubeID)
	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	logrus.Debugf("Get cloud profile %s", k.ProfileID)
	kubeProfile, err := h.profileSvc.Get(r.Context(), k.ProfileID)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, k.ProfileID, err)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	config, err := steps.NewConfigFromKube(kubeProfile, k)
	if err != nil {
		logrus.Errorf("New config %v", err.Error())
		message.SendUnknownError(w, err)
		return
	}

	logrus.Debugf("load clout specific data from kube %s", k.ID)
	// Load things specific to cloud provider
	err = util.LoadCloudSpecificDataFromKube(k, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	logrus.Debugf("Get cloud account %s", k.AccountName)
	acc, err := h.accountService.Get(r.Context(), k.AccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	logrus.Debug("Fill config with cloud account credentials")
	err = util.FillCloudAccountCredentials(acc, config)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	logrus.Debugf("Restart cluster %s provisioning", k.ID)
	err = h.kubeProvisioner.RestartClusterProvisioning(r.Context(),
		kubeProfile, config, k.Tasks)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) importKube(w http.ResponseWriter, r *http.Request) {
	type importRequest struct {
		KubeConfig            string            `json:"kubeconfig"`
		ClusterName           string            `json:"clusterName"`
		CloudAccountName      string            `json:"cloudAccountName"`
		PublicKey             string            `json:"publicKey"`
		PrivateKey            string            `json:"privateKey"`
		CloudSpecificSettings map[string]string `json:"cloudSpecificSettings"`
	}

	var req importRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		logrus.Error(err)
		message.SendInvalidJSON(w, err)
		return
	}

	kubeConfig, err := clientcmd.Load([]byte(req.KubeConfig))

	if err != nil {
		logrus.Error(err)
		message.SendInvalidJSON(w, err)
		return
	}

	config := &steps.Config{
		CloudAccountName: req.CloudAccountName,
		ClusterName:      req.ClusterName,
		IsBootstrap:      true,
		Kube: model.Kube{
			SSHConfig: model.SSHConfig{
				Port:                "22",
				Timeout:             10,
				BootstrapPrivateKey: req.PrivateKey,
				BootstrapPublicKey:  req.PublicKey,
			},
		},
		Masters: steps.NewMap(map[string]*model.Machine{}),
		Nodes:   steps.NewMap(map[string]*model.Machine{}),
	}

	importTask, err := workflows.NewTask(config, workflows.ImportCluster, h.repo)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	kube, err := kubeFromKubeConfig(*kubeConfig)

	if err != nil {
		message.SendInvalidCredentials(w, err)
		return
	}
	// Grab all k8s nodes from kube-apiserver
	nodes, err := h.svc.ListNodes(r.Context(), kube, "")
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	cloudAccount, err := h.accountService.Get(r.Context(), req.CloudAccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, req.CloudAccountName, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	err = util.FillCloudAccountCredentials(cloudAccount, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	var clusterID string

	if len(importTask.ID) > 8 {
		clusterID = importTask.ID[:8]
	} else {
		message.SendValidationFailed(w, errors.New("import task id is too short"))
		return
	}

	config.ClusterID = clusterID
	config.AWSConfig.Region = "ap-south-1"

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(struct {
		ClusterID string `json:"clusterId"`
	}{
		ClusterID: clusterID,
	})

	if err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	go func() {
		fileName := util.MakeFileName(importTask.ID)
		writer, err := h.getWriter(fileName)

		if err != nil {
			message.SendUnknownError(w, err)
			return
		}

		for _, node := range nodes {
			machine := model.Machine{}

			for _, address := range node.Status.Addresses {
				if address.Type == "ExternalIP" {
					machine.PublicIp = address.Address
				} else if address.Type == "InternalIP" {
					machine.PrivateIp = address.Address
				}
			}

			machine.Role = model.RoleNode

			for key := range node.Labels {
				if strings.HasSuffix(key, "master") {
					machine.Role = model.RoleMaster
				}
			}

			if machine.Role == model.RoleMaster {
				config.AddMaster(&machine)
			} else {
				config.AddNode(&machine)
			}
		}

		importTask.Config = config
		resultChan := importTask.Run(context.Background(), *importTask.Config, writer)
		err = <-resultChan

		if err != nil {
			logrus.Errorf("task %s has finished with error %v", importTask.ID, err)
		}

		if err := createKubeFromConfig(context.Background(), config, h.svc); err != nil {
			logrus.Errorf("Error creating the cluster")
		}

		logrus.Infof("Import task %s has successfully finished", importTask.ID)
	}()
}
