package controlplane

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
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
	sshRunner "github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/sghelm"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/testutils/assert"
	"github.com/supergiant/control/pkg/user"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/etcd"
	"github.com/supergiant/control/pkg/workflows/steps/flannel"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
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

// Config is the server configuration
type Config struct {
	Port          int
	Addr          string
	EtcdUrl       string
	TemplatesDir  string
	SpawnInterval time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
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
	if err := generateUserIfColdStart(cfg); err != nil {
		return nil, err
	}

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

//generateUserIfColdStart checks if there are any users in the db and if not (i.e. on first launch) generates a root user
func generateUserIfColdStart(cfg *Config) error {
	etcdCfg := clientv3.Config{
		DialTimeout: time.Second * 10,
		Endpoints:   []string{cfg.EtcdUrl},
	}
	repository := storage.NewETCDRepository(etcdCfg)
	userService := user.NewService(user.DefaultStoragePrefix, repository)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := userService.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		u := &user.User{
			Login:    "root",
			Password: util.RandomString(13),
		}
		logrus.Infof("first time launch detected, use %s as login and %s as password", u.Login, u.Password)
		err := userService.Create(ctx, u)
		if err != nil {
			return nil
		}
	}

	return nil
}

func validate(cfg *Config) error {
	if cfg.EtcdUrl == "" {
		return errors.New("etcd url can't be empty")
	}

	if err := assert.CheckETCD(cfg.EtcdUrl); err != nil {
		return errors.Wrapf(err, "etcd url %s", cfg.EtcdUrl)
	}

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
	etcdCfg := clientv3.Config{
		Endpoints: []string{cfg.EtcdUrl},
	}
	router := mux.NewRouter()

	protectedAPI := router.PathPrefix("/v1/api").Subrouter()
	repository := storage.NewETCDRepository(etcdCfg)

	accountService := account.NewService(account.DefaultStoragePrefix, repository)
	accountHandler := account.NewHandler(accountService)
	accountHandler.Register(protectedAPI)

	//TODO Add generation of jwt token
	jwtService := jwt.NewTokenService(86400, []byte("test"))
	userService := user.NewService(user.DefaultStoragePrefix, repository)
	userHandler := user.NewHandler(userService, jwtService)

	router.HandleFunc("/auth", userHandler.Authenticate).Methods(http.MethodPost)
	//Opening it up for testing right now, will be protected after implementing initial user generation
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
	flannel.Init()
	kubelet.Init()
	manifest.Init()
	poststart.Init()
	tiller.Init()
	etcd.Init()
	ssh.Init()
	network.Init()
	clustercheck.Init()
	prometheus.Init()
	gce.Init()

	amazon.InitImportKeyPair(amazon.GetEC2)
	amazon.InitCreateMachine(amazon.GetEC2)
	amazon.InitCreateSecurityGroups(amazon.GetEC2)
	amazon.InitCreateVPC(amazon.GetEC2)
	amazon.InitCreateSubnet(amazon.GetEC2)
	amazon.InitDeleteCluster(amazon.GetEC2)
	amazon.InitDeleteNode(amazon.GetEC2)

	workflows.Init()

	taskHandler := workflows.NewTaskHandler(repository, sshRunner.NewRunner, accountService)
	taskHandler.Register(router)

	helmService, err := sghelm.NewService(repository)
	if err != nil {
		return nil, errors.Wrap(err, "new helm service")
	}
	go ensureHelmRepositories(helmService)

	helmHandler := sghelm.NewHandler(helmService)
	helmHandler.Register(protectedAPI)

	kubeService := kube.NewService(kube.DefaultStoragePrefix,
		repository, helmService)

	taskProvisioner := provisioner.NewProvisioner(repository,
		kubeService,
		cfg.SpawnInterval)
	tokenGetter := provisioner.NewEtcdTokenGetter()
	provisionHandler := provisioner.NewHandler(kubeService, accountService,
		tokenGetter, taskProvisioner)
	provisionHandler.Register(protectedAPI)

	kubeHandler := kube.NewHandler(kubeService, accountService,
		taskProvisioner, repository)
	kubeHandler.Register(protectedAPI)

	authMiddleware := api.Middleware{
		TokenService: jwtService,
	}
	protectedAPI.Use(authMiddleware.AuthMiddleware, api.ContentTypeJSON)

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
