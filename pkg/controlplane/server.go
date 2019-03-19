package controlplane

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/api"
	"github.com/supergiant/control/pkg/jwt"
	"github.com/supergiant/control/pkg/kube"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/provisioner"
	"github.com/supergiant/control/pkg/proxy"
	sshRunner "github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/sghelm"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/user"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/drain"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
	"github.com/supergiant/control/pkg/workflows/steps/storageclass"
	"github.com/supergiant/control/pkg/workflows/steps/tiller"
)

type Server struct {
	server http.Server
	cfg    *Config
}

func (srv *Server) Start() {
	err := srv.server.ListenAndServe()
	if err != nil {
		logrus.Error(err)
	}
}

func (srv *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	err := srv.server.Shutdown(ctx)

	if err != nil {
		logrus.Error(err)
	}
}

func New(cfg *Config) (*Server, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}

	r, err := configureApplication(cfg)
	if err != nil {
		return nil, err
	}

	s := NewServer(r, cfg)

	return s, nil
}

func NewServer(router *mux.Router, cfg *Config) *Server {
	headersOk := handlers.AllowedHeaders([]string{
		"Access-Control-Request-Headers",
		"Authorization",
	})
	methodsOk := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodOptions,
		http.MethodDelete,
	})

	// TODO add TLS support
	s := &Server{
		cfg: cfg,
		server: http.Server{
			Handler:      handlers.CORS(headersOk, methodsOk)(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router)),
			Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
	http.DefaultClient.Timeout = cfg.IdleTimeout

	return s
}

func validate(cfg *Config) error {
	if cfg.Port <= 0 {
		return errors.New("port can't be negative")
	}

	if cfg.SpawnInterval == 0 {
		return errors.New("spawn interval must not be 0")
	}

	return nil
}

func configureApplication(cfg *Config) (*mux.Router, error) {
	//TODO will work for now, but we should revisit ETCD configuration later
	router := mux.NewRouter()

	protectedAPI := router.PathPrefix("/v1/api").Subrouter()
	repository, err := storage.GetStorage(cfg.StorageMode, cfg.StorageURI)

	if err != nil {
		return nil, errors.Wrapf(err, "get storage type %s uri %s",
			cfg.StorageMode, cfg.StorageURI)
	}

	accountService := account.NewService(account.DefaultStoragePrefix, repository)
	accountHandler := account.NewHandler(accountService)
	accountHandler.Register(protectedAPI)

	//TODO Add generation of jwt token
	jwtService := jwt.NewTokenService(86400, []byte("test"))
	userService := user.NewService(user.DefaultStoragePrefix, repository)
	userHandler := user.NewHandler(userService, jwtService)

	router.HandleFunc("/version", NewVersionHandler(cfg.Version))
	router.HandleFunc("/auth", userHandler.Authenticate).Methods(http.MethodPost)
	router.HandleFunc("/root", userHandler.RegisterRootUser).Methods(http.MethodPost)
	router.HandleFunc("/coldstart", userHandler.IsColdStart).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost)

	profileService := profile.NewService(profile.DefaultKubeProfilePreifx, repository)
	kubeProfileHandler := profile.NewHandler(profileService)
	kubeProfileHandler.Register(protectedAPI)

	// Read templates first and then initialize workflows with steps that uses these templates
	if err := templatemanager.Init(cfg.TemplatesDir); err != nil {
		return nil, err
	}

	digitalocean.Init()
	certificates.Init()
	authorizedKeys.Init()
	cni.Init()
	docker.Init()
	downloadk8sbinary.Init()
	kubelet.Init()
	poststart.Init()
	tiller.Init()
	ssh.Init()
	network.Init()
	clustercheck.Init()
	prometheus.Init()
	gce.Init()
	storageclass.Init()
	drain.Init()
	kubeadm.Init()

	amazon.InitFindAMI(amazon.GetEC2)
	amazon.InitImportKeyPair(amazon.GetEC2)
	amazon.InitCreateInstanceProfiles(amazon.GetIAM)
	amazon.InitCreateMachine(amazon.GetEC2)
	amazon.InitCreateSecurityGroups(amazon.GetEC2)
	amazon.InitCreateVPC(amazon.GetEC2)
	amazon.InitCreateSubnet(amazon.GetEC2, accountService)
	amazon.InitDeleteClusterMachines(amazon.GetEC2)
	amazon.InitDeleteNode(amazon.GetEC2)
	amazon.InitDeleteSecurityGroup(amazon.GetEC2)
	amazon.InitDeleteVPC(amazon.GetEC2)
	amazon.InitDeleteSubnets(amazon.GetEC2)
	amazon.InitCreateRouteTable(amazon.GetEC2)
	amazon.InitAssociateRouteTable(amazon.GetEC2)
	amazon.InitCreateInternetGateway(amazon.GetEC2)
	amazon.InitDeleteSubnets(amazon.GetEC2)
	amazon.InitDisassociateRouteTable(amazon.GetEC2)
	amazon.InitDeleteRouteTable(amazon.GetEC2)
	amazon.InitDeleteInternetGateWay(amazon.GetEC2)
	amazon.InitDeleteKeyPair(amazon.GetEC2)
	workflows.Init()

	taskHandler := workflows.NewTaskHandler(repository, sshRunner.NewRunner, accountService)
	taskHandler.Register(protectedAPI)

	helmService, err := sghelm.NewService(repository)
	if err != nil {
		return nil, errors.Wrap(err, "new helm service")
	}
	if coldstart, err := userService.IsColdStart(context.Background()); err == nil && coldstart {
		go ensureHelmRepositories(helmService)
	} else if err != nil {
		return nil, err
	}

	helmHandler := sghelm.NewHandler(helmService)
	helmHandler.Register(protectedAPI)

	kubeService := kube.NewService(kube.DefaultStoragePrefix,
		repository, helmService)

	taskProvisioner := provisioner.NewProvisioner(repository,
		kubeService,
		cfg.SpawnInterval)
	provisionHandler := provisioner.NewHandler(kubeService, accountService,
		profileService, taskProvisioner)
	provisionHandler.Register(protectedAPI)
	apiProxy := proxy.NewReverseProxyContainer(cfg.ProxiesPortRange,
		logrus.New().WithField("component", "proxy"))

	kubeHandler := kube.NewHandler(kubeService, accountService,
		profileService, taskProvisioner, taskProvisioner,
		repository, apiProxy)
	kubeHandler.Register(protectedAPI)

	authMiddleware := api.Middleware{
		TokenService: jwtService,
	}
	protectedAPI.Use(authMiddleware.AuthMiddleware, api.ContentTypeJSON)

	if cfg.PprofListenStr != "" {
		go func() {
			logrus.Debugf("Start pprof on %s", cfg.PprofListenStr)
			logrus.Info(http.ListenAndServe(cfg.PprofListenStr, nil))
		}()
	}

	if err := ServeUI(cfg, router); err != nil {
		return nil, err
	}
	return router, nil
}

func ensureHelmRepositories(svc sghelm.Servicer) {
	if svc == nil {
		return
	}

	entries := []repo.Entry{
		{
			Name: "supergiant",
			URL:  "https://supergiant.github.io/charts",
		},
		{
			Name: "stable",
			URL:  "https://kubernetes-charts.storage.googleapis.com",
		},
	}

	for _, entry := range entries {
		_, err := svc.CreateRepo(context.Background(), &entry)
		if err != nil {
			if !sgerrors.IsAlreadyExists(err) {
				logrus.Errorf("failed to add %q helm repository: %v", entry.Name, err)
			}
			continue
		}
		logrus.Infof("helm repository has been added: %s", entry.Name)
	}

}

func NewVersionHandler(version string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	}
}
