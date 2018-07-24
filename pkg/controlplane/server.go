package controlplane

import (
	"net/http"

	"fmt"
	"time"

	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/api"
	"github.com/supergiant/supergiant/pkg/helm"
	"github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/kube"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
	"github.com/supergiant/supergiant/pkg/user"
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
	err := srv.server.Close()
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

	// TODO add TLS support
	s := &Server{
		cfg: cfg,
		server: http.Server{
			Handler:      handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(r),
			Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 15,
			IdleTimeout:  time.Second * 120,
		},
	}

	return s, nil
}

func validate(cfg *Config) error {
	if cfg.EtcdUrl == "" {
		return errors.New("etcd url can't be empty")
	}

	if err := assert.CheckETCD(cfg.EtcdUrl); err != nil {
		return err
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

	kubeService := kube.NewService(kube.DefaultStoragePrefix, repository)
	kubeHandler := kube.NewHandler(kubeService)
	kubeHandler.Register(protectedAPI)

	//TODO Add generation of jwt token
	jwtService := jwt.NewTokenService(86400, []byte("test"))
	userService := user.NewService(user.DefaultStoragePrefix, repository)
	userHandler := user.NewHandler(userService, jwtService)

	router.HandleFunc("/auth", userHandler.Authenticate).Methods(http.MethodPost)
	//Opening it up for testing right now, will be protected after implementing initial user generation
	router.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost)

	kubeProfileService := profile.NewKubeProfileService(profile.DefaultKubeProfilePreifx, repository)
	kubeProfileHandler := profile.NewKubeProfileHandler(kubeProfileService)
	kubeProfileHandler.Register(protectedAPI)

	nodeProfileService := profile.NewNodeProfileService(profile.DefaultNodeProfilePrefix, repository)
	nodeProfileHandler := profile.NewNodeProfileHandler(nodeProfileService)
	nodeProfileHandler.Register(protectedAPI)

	helmService := helm.NewService(repository)
	helmHandler := helm.NewHandler(helmService)
	helmHandler.Register(protectedAPI)

	authMiddleware := api.Middleware{
		TokenService: jwtService,
		UserService:  userService,
	}
	protectedAPI.Use(authMiddleware.AuthMiddleware)

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
