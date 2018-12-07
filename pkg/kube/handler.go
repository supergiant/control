package kube

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/proxy"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	cache  metricCache
	m      sync.Mutex
	client *http.Client
)

type metricCache struct {
	m    sync.RWMutex
	data map[string]entry
}

func (c *metricCache) get(key string) *MetricResponse {
	c.m.RLock()
	defer c.m.RUnlock()
	e := c.data[key]

	if e.timestamp == 0 {
		return nil
	}

	if e.timestamp < time.Now().Unix()-60 {
		return nil
	}

	return e.value
}

func (c *metricCache) set(key string, value *MetricResponse) {
	c.m.Lock()
	defer c.m.Unlock()
	c.data[key] = entry{
		timestamp: time.Now().Unix(),
		value:     value,
	}
}

type entry struct {
	timestamp int64
	value     *MetricResponse
}

type accountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type nodeProvisioner interface {
	ProvisionNodes(context.Context, []profile.NodeProfile, *model.Kube,
		*steps.Config) ([]string, error)
	// Method that cancels newly added nodes to working cluster
	Cancel(string) error
}


type kubeProvisioner interface{
	RestartClusterProvisioning(context.Context, profile.Profile, *steps.Config, map[string][]string) error
}

// TODO(stgleb): use standard k8s structs for that
type k8SServices struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string            `json:"name"`
			Namespace         string            `json:"namespace"`
			SelfLink          string            `json:"selfLink"`
			UID               string            `json:"uid"`
			ResourceVersion   string            `json:"resourceVersion"`
			CreationTimestamp time.Time         `json:"creationTimestamp"`
			Labels            map[string]string `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			Ports []struct {
				Name     string `json:"name"`
				Protocol string `json:"protocol"`
				Port     int    `json:"port"`
			} `json:"ports"`
			Selector struct {
				App string `json:"app"`
			} `json:"selector"`
			ClusterIP       string `json:"clusterIP"`
			Type            string `json:"type"`
			SessionAffinity string `json:"sessionAffinity"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct {
			} `json:"loadBalancer"`
		} `json:"status"`
	} `json:"items"`
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
	workflowMap     map[clouds.Name]workflows.WorkflowSet
	repo            storage.Interface
	getWriter       func(string) (io.WriteCloser, error)
	getMetrics      func(string, *model.Kube) (*MetricResponse, error)
	proxies         proxy.Container
}

func init() {
	cache = metricCache{
		data: make(map[string]entry),
	}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: time.Second * 30,
		MaxIdleConnsPerHost: 100,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
}

// Use shared http.Client to perform request
func doReq(req *http.Request) (*http.Response, error) {
	m.Lock()
	defer m.Unlock()
	return client.Do(req)
}

// NewHandler constructs a Handler for kubes.
func NewHandler(
	svc Interface,
	accountService accountGetter,
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
		workflowMap: map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				DeleteCluster: workflows.DigitalOceanDeleteCluster,
				DeleteNode:    workflows.DigitalOceanDeleteNode,
			},
			clouds.AWS: {
				DeleteCluster: workflows.AWSDeleteCluster,
				DeleteNode:    workflows.AWSDeleteNode,
			},
			clouds.GCE: {
				DeleteCluster: workflows.GCEDeleteCluster,
				DeleteNode:    workflows.GCEDeleteNode,
			},
		},
		repo:      repo,
		getWriter: util.GetWriter,
		getMetrics: func(metricURI string, k *model.Kube) (*MetricResponse, error) {
			if m := cache.get(metricURI); m != nil {
				logrus.Debugf("metric cache hit")
				return m, nil
			}
			logrus.Debugf("metric cache miss")

			metricResponse := &MetricResponse{}
			// TODO(stgleb): Add caching for metric
			req, err := http.NewRequest(http.MethodGet, metricURI, nil)

			if err != nil {
				return nil, err
			}

			req.SetBasicAuth(k.User, k.Password)
			logrus.Debugf("Get metric with URI %s for kube %s", metricURI, k.ID)
			resp, err := doReq(req)

			if err != nil {
				return nil, err
			}

			err = json.NewDecoder(resp.Body).Decode(metricResponse)

			if err != nil {
				return nil, err
			}

			resp.Body.Close()
			cache.set(metricURI, metricResponse)

			return metricResponse, nil
		},
		proxies: proxies,
	}
}

// Register adds kube handlers to a router.
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/kubes", h.createKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes", h.listKubes).Methods(http.MethodGet)
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

	r.HandleFunc("/kubes/{kubeID}/nodes", h.addNode).Methods(http.MethodPost)
	r.HandleFunc("/kubes/{kubeID}/nodes/{nodename}", h.deleteNode).Methods(http.MethodDelete)
	r.HandleFunc("/kubes/{kubeID}/metrics", h.getClusterMetrics).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{kubeID}/nodes/metrics", h.getNodesMetrics).Methods(http.MethodGet)
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

	t, err := workflows.NewTask(h.workflowMap[acc.Provider].DeleteCluster, h.repo)

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

	// Load things specific to cloud provider
	err = util.LoadCloudSpecificDataFromKube(k, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	err = util.FillCloudAccountCredentials(r.Context(), acc, config)

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

		// Finally delete cluster record from etcd
		if err := h.svc.Delete(context.Background(), kubeID); err != nil {
			logrus.Errorf("delete kube %s caused %v", kubeID, err)
			return
		}

		h.deleteClusterTasks(context.Background(), kubeID)
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

// Add node to working kube
func (h *Handler) addNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeID := vars["kubeID"]
	k, err := h.svc.Get(r.Context(), kubeID)

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
		Zone:            k.Zone,
		Arch:            k.Arch,
		OperatingSystem: k.OperatingSystem,
		UbuntuVersion:   k.OperatingSystemVersion,
		DockerVersion:   k.DockerVersion,
		K8SVersion:      k.K8SVersion,
		HelmVersion:     k.HelmVersion,
		User:            k.User,
		Password:        k.Password,

		NetworkType:           k.Networking.Type,
		CIDR:                  k.Networking.CIDR,
		FlannelVersion:        k.Networking.Version,
		CloudSpecificSettings: k.CloudSpec,

		NodesProfiles: []profile.NodeProfile{
			{},
		},

		RBACEnabled: k.RBACEnabled,
	}

	config := steps.NewConfig(k.Name, "", k.AccountName, kubeProfile)
	config.ClusterID = k.ID
	config.CertificatesConfig.CAKey = k.Auth.CAKey
	config.CertificatesConfig.CACert = k.Auth.CACert
	config.CertificatesConfig.AdminCert = k.Auth.AdminCert
	config.CertificatesConfig.AdminKey = k.Auth.AdminKey

	if len(k.Masters) != 0 {
		config.AddMaster(util.GetRandomNode(k.Masters))
	} else {
		http.Error(w, "no master found", http.StatusNotFound)
		return
	}

	// Get cloud account fill appropriate config structure
	// with cloud account credentials
	err = util.FillCloudAccountCredentials(r.Context(), acc, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)
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
func (h *Handler) deleteNode(w http.ResponseWriter, r *http.Request) {
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

	var n *node.Node

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

	t, err := workflows.NewTask(h.workflowMap[acc.Provider].DeleteNode, h.repo)

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
		Node: *n,
	}

	err = util.FillCloudAccountCredentials(r.Context(), acc, config)

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

	util.MakeFileName(t.ID)
	writer, err := h.getWriter(t.ID)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	errChan := t.Run(context.Background(), *config, writer)

	// Update cluster state when deletion completes
	go func() {
		// Set node to deleting state
		nodeToDelete, ok := k.Nodes[nodeName]

		if !ok {
			logrus.Errorf("Node %s not found", nodeName)
			return
		}
		nodeToDelete.State = node.StateDeleting
		k.Nodes[nodeName] = nodeToDelete
		err := h.svc.Create(context.Background(), k)

		if err != nil {
			logrus.Errorf("update cluster %s caused %v", kubeID, err)
		}

		err = <-errChan

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
		masterNode *node.Node
		response   = map[string]interface{}{}
		baseUrl    = "api/v1/namespaces/default/services/prometheus-operated:9090/proxy"
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
		masterNode *node.Node
		response   = map[string]map[string]interface{}{}
		baseUrl    = "api/v1/namespaces/default/services/prometheus-operated:9090/proxy"
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
	var (
		servicesUrl = "api/v1/services"
		kubeID      string
		masterNode  *node.Node
		k8sServices = &k8SServices{}
	)
	vars := mux.Vars(r)
	kubeID = vars["kubeID"]

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

	serviceURL := fmt.Sprintf("https://%s/%s", masterNode.PublicIp, servicesUrl)
	req, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	req.SetBasicAuth(k.User, k.Password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(k8sServices)

	if err != nil {
		logrus.Error(err)
		message.SendInvalidJSON(w, err)
		return
	}

	// TODO(stgleb): Figure out which ports are worth to be proxy
	webPorts := map[string]struct{}{
		"web":     {},
		"http":    {},
		"https":   {},
		"service": {},
	}

	var serviceInfos = make([]*ServiceInfo, 0)
	var targetServices = make([]*proxy.Target, 0)

	for _, service := range k8sServices.Items {
		if !contains("kubernetes.io/cluster-service", "true", service.Metadata.Labels) {
			continue
		}

		if len(service.Spec.Ports) == 1 && service.Spec.Ports[0].Protocol == "TCP" {
			var serviceInfo = &ServiceInfo{
				ID:        service.Metadata.UID,
				Name:      service.Metadata.Name,
				Type:      service.Spec.Type,
				Namespace: service.Metadata.Namespace,
			}
			serviceInfos = append(serviceInfos, serviceInfo)

			targetServices = append(targetServices, &proxy.Target{
				ProxyID: kubeID + service.Metadata.UID,
				TargetURL: fmt.Sprintf("https://%s%s:%d/proxy",
					masterNode.PublicIp, service.Metadata.SelfLink, service.Spec.Ports[0].Port),
				SelfLink: service.Metadata.SelfLink,
				User:     k.User,
				Password: k.Password,
			})
			continue
		}

		for _, port := range service.Spec.Ports {
			if port.Protocol == "TCP" {
				if _, ok := webPorts[port.Name]; ok {
					var serviceInfo = &ServiceInfo{
						ID:        service.Metadata.UID,
						Name:      service.Metadata.Name,
						Type:      service.Spec.Type,
						Namespace: service.Metadata.Namespace,
					}
					serviceInfos = append(serviceInfos, serviceInfo)

					targetServices = append(targetServices, &proxy.Target{
						ProxyID: kubeID + service.Metadata.UID,
						TargetURL: fmt.Sprintf("https://%s%s:%d/proxy",
							masterNode.PublicIp, service.Metadata.SelfLink, port.Port),
						SelfLink: service.Metadata.SelfLink,
						User:     k.User,
						Password: k.Password,
					})
				}
			}
		}
	}

	// TODO implement proxy removing logic under separate ticket
	err = h.proxies.RegisterProxies(targetServices)
	if err != nil {
		logrus.Error(err)
		message.SendUnknownError(w, err)
		return
	}

	var proxies = h.proxies.GetProxies(kubeID)

	for _, service := range serviceInfos {
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

	k, err := h.svc.Get(r.Context(), kubeID)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kubeID, err)
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

	// Load things specific to cloud provider
	err = util.LoadCloudSpecificDataFromKube(k, config)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	err = h.kubeProvisioner.RestartClusterProvisioning(r.Context(), profile.Profile{}, config, k.Tasks)

	if err != nil {
		message.SendUnknownError(w, err)
		return
	}
}