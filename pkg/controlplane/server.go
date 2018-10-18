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

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/api"
	"github.com/supergiant/supergiant/pkg/helm"
	"github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/kube"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/provisioner"
	sshRunner "github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
	"github.com/supergiant/supergiant/pkg/user"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps/amazon"
	"github.com/supergiant/supergiant/pkg/workflows/steps/certificates"
	"github.com/supergiant/supergiant/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/supergiant/pkg/workflows/steps/cni"
	"github.com/supergiant/supergiant/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
	"github.com/supergiant/supergiant/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
	"github.com/supergiant/supergiant/pkg/workflows/steps/flannel"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubelet"
	"github.com/supergiant/supergiant/pkg/workflows/steps/manifest"
	"github.com/supergiant/supergiant/pkg/workflows/steps/network"
	"github.com/supergiant/supergiant/pkg/workflows/steps/poststart"
	"github.com/supergiant/supergiant/pkg/workflows/steps/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps/tiller"
)

type Server struct {
	server http.Server
	cfg    *Config
}

func (srv *Server) Start() {
	err := srv.server.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}

func (srv *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	err := srv.server.Shutdown(ctx)

	if err != nil {
		logrus.Fatal(err)
	}
}

// Config is the server configuration
type Config struct {
	Port         int
	Addr         string
	EtcdUrl      string
	LogLevel     string
	TemplatesDir string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(cfg *Config) (*Server, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}

	configureLogging(cfg)
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

	profileService := profile.NewKubeProfileService(profile.DefaultKubeProfilePreifx, repository)
	kubeProfileHandler := profile.NewKubeProfileHandler(profileService)
	kubeProfileHandler.Register(protectedAPI)

	// Read templates first and then initialize workflows with steps that uses these templates
	if err := templatemanager.Init(cfg.TemplatesDir); err != nil {
		return nil, err
	}

	digitalocean.Init()
	certificates.Init()
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
	amazon.InitImportKeyPair(amazon.GetEC2)
	amazon.InitCreateMachine(amazon.GetEC2)
	amazon.InitCreateSecurityGroups(amazon.GetEC2)
	amazon.InitCreateVPC(amazon.GetEC2)
	amazon.InitCreateSubnet(amazon.GetEC2)

	workflows.Init()

	taskHandler := workflows.NewTaskHandler(repository, sshRunner.NewRunner, accountService)
	taskHandler.Register(router)

	kubeService := kube.NewService(kube.DefaultStoragePrefix, repository)

	taskProvisioner := provisioner.NewProvisioner(repository, kubeService)
	tokenGetter := provisioner.NewEtcdTokenGetter()
	provisionHandler := provisioner.NewHandler(accountService, tokenGetter, taskProvisioner)
	provisionHandler.Register(protectedAPI)

	kubeHandler := kube.NewHandler(kubeService, accountService, taskProvisioner, repository)
	kubeHandler.Register(protectedAPI)

	helmService, err := helm.NewService(repository)
	if err != nil {
		return nil, errors.Wrap(err, "new helm service")
	}
	go ensureHelmRepositories(helmService)

	helmHandler := helm.NewHandler(helmService)
	helmHandler.Register(protectedAPI)

	authMiddleware := api.Middleware{
		TokenService: jwtService,
		UserService:  userService,
	}
	protectedAPI.Use(authMiddleware.AuthMiddleware, api.ContentTypeJSON)

	return router, nil
}

func configureLogging(cfg *Config) {
	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Warnf("incorrect logging level %s provided, setting INFO as default...", l)
		logrus.SetLevel(logrus.InfoLevel)
		return
	}
	logrus.SetLevel(l)
}

func ensureHelmRepositories(svc helm.Servicer) {
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
