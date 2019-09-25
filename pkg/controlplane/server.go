package controlplane

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/supergiant/control/pkg/workflows/steps/helm"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
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
	"github.com/supergiant/control/pkg/workflows/steps/apply"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedkeys"
	"github.com/supergiant/control/pkg/workflows/steps/azure"
	"github.com/supergiant/control/pkg/workflows/steps/bootstraptoken"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/cloudcontroller"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/configmap"
	"github.com/supergiant/control/pkg/workflows/steps/dashboard"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/drain"
	"github.com/supergiant/control/pkg/workflows/steps/evacuate"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
	"github.com/supergiant/control/pkg/workflows/steps/install_app"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
	"github.com/supergiant/control/pkg/workflows/steps/storageclass"
	"github.com/supergiant/control/pkg/workflows/steps/tiller"
	"github.com/supergiant/control/pkg/workflows/steps/uncordon"
	"github.com/supergiant/control/pkg/workflows/steps/upgrade"
	_ "github.com/supergiant/control/statik"
)

type Server struct {
	server http.Server
	cfg    *Config
}

func (srv *Server) Start() {
	logrus.Infof("configuratino: %+v", srv.cfg)
	logrus.Infof("supergiant is listening on %s", srv.server.Addr)

	var err error
	if srv.server.TLSConfig != nil {
		err = srv.server.ListenAndServeTLS(srv.cfg.CertFile, srv.cfg.KeyFile)
	} else {
		err = srv.server.ListenAndServe()
	}
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

// Config is the server configuration
type Config struct {
	Port         int
	InsecurePort int
	CertFile     string
	KeyFile      string
	Addr         string
	StorageMode  string
	StorageURI   string
	TemplatesDir string
	LogDir       string

	SpawnInterval time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	PprofListenStr string

	ProxiesPortRange proxy.PortRange

	Version string
}

func New(cfg *Config) (*Server, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}

	r, err := configureApplication(cfg)
	if err != nil {
		return nil, err
	}

	return NewServer(r, cfg)
}

func NewServer(router *mux.Router, cfg *Config) (*Server, error) {
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

	port := cfg.InsecurePort
	var tlsCfg *tls.Config
	if cfg.Port != 0 {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "load server certificates")
		}

		port = cfg.Port
		tlsCfg = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	return &Server{
		cfg: cfg,
		server: http.Server{
			Handler:      handlers.CORS(headersOk, methodsOk)(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router)),
			Addr:         fmt.Sprintf("%s:%d", cfg.Addr, port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			TLSConfig:    tlsCfg,
		},
	}, nil
}

func validate(cfg *Config) error {
	if cfg.SpawnInterval == 0 {
		return errors.New("spawn interval must not be 0")
	}

	return nil
}

func configureApplication(cfg *Config) (*mux.Router, error) {
	//TODO will work for now, but we should revisit ETCD configuration later
	router := mux.NewRouter()

	protectedAPI := router.PathPrefix("/api/v1").Subrouter()
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
		return nil, errors.Wrap(err, "templatemanager: init")
	}

	digitalocean.Init()
	certificates.Init()
	authorizedkeys.Init()
	cni.Init()
	docker.Init()
	downloadk8sbinary.Init()
	kubelet.Init()
	poststart.Init()
	tiller.Init()
	ssh.Init()
	network.Init()
	clustercheck.Init()
	cloudcontroller.Init()
	prometheus.Init()
	dashboard.Init()
	gce.Init(accountService)
	storageclass.Init()
	drain.Init()
	kubeadm.Init()
	bootstraptoken.Init()
	configmap.Init()
	upgrade.Init()
	uncordon.Init()
	evacuate.Init()
	install_app.Init()
	helm.Init()

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
	amazon.InitCreateLoadBalancer(amazon.GetELB)
	amazon.InitDeleteLoadBalancer(amazon.GetELB)
	amazon.InitRegisterInstance(amazon.GetELB)
	amazon.InitImportClusterStep(amazon.GetEC2)
	amazon.InitImportSubnetDescriber(amazon.GetEC2)
	amazon.InitImportInternetGatewayStep(amazon.GetEC2)
	amazon.InitImportRouteTablesStep(amazon.GetEC2)
	amazon.InitCreateTagsStep(amazon.GetEC2)
	apply.Init()
	azure.Init()

	workflows.Init()

	taskHandler := workflows.NewTaskHandler(repository, sshRunner.NewRunner, accountService, cfg.LogDir)
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
		cfg.SpawnInterval, cfg.LogDir)
	provisionHandler := provisioner.NewHandler(kubeService, accountService,
		profileService, taskProvisioner)
	provisionHandler.Register(protectedAPI)
	apiProxy := proxy.NewReverseProxyContainer(cfg.ProxiesPortRange,
		logrus.New().WithField("component", "proxy"))

	kubeHandler := kube.NewHandler(kubeService, accountService,
		profileService, taskProvisioner, taskProvisioner, helmService,
		repository, apiProxy, cfg.LogDir)
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

	if err := serveUI(cfg, router); err != nil {
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

func serveUI(cfg *Config, router *mux.Router) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(trimPrefix(http.FileServer(statikFS)))
	return nil
}

func trimPrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This code path is for static resources
		if len(strings.Split(r.URL.Path, ".")) == 1 {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = "/"
			logrus.Debugf("Change path URL %s to %s",
				r.URL.Path, r2.URL.Path)
			h.ServeHTTP(w, r2)
		} else {
			// This codepath is for URL paths from browser line
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = r2.URL.Path[1:]
			logrus.Debugf("Change asset URL %s to %s",
				r.URL.Path, r2.URL.Path)
			h.ServeHTTP(w, r2)
		}
	})
}

func NewVersionHandler(version string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	}
}
