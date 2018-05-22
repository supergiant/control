package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kclient "k8s.io/client-go/kubernetes"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

type NodeSize struct {
	Name     string  `json:"name"`
	RAMGIB   float64 `json:"ram_gib"`
	CPUCores float64 `json:"cpu_cores"`
}

type Settings struct {
	ConfigFilePath string

	PsqlHost               string `json:"psql_host"`
	PsqlDb                 string `json:"psql_db"`
	PsqlUser               string `json:"psql_user"`
	PsqlPass               string `json:"psql_pass"`
	SQLiteFile             string `json:"sqlite_file"`
	PublishHost            string `json:"publish_host" default:"localhost"`
	HTTPPort               string `json:"http_port" default:"8080"`
	HTTPSPort              string `json:"https_port"`
	SSLCertFile            string `json:"ssl_cert_file"`
	SSLKeyFile             string `json:"ssl_key_file"`
	LogPath                string `json:"log_file"`
	LogLevel               string `json:"log_level"`
	SupportPassword        string `json:"support_password"`
	UIEnabled              bool   `json:"ui_enabled"`
	CapacityServiceEnabled bool   `json:"capacity_service_enabled"`

	// NOTE these MUST be provided in ascending order by cost in order to
	// correctly provision the smallest size on Kube creation
	//
	// NodeSizes is a map of provider name (ex. "aws") and node sizes
	NodeSizes map[string][]*NodeSize `json:"node_sizes"`

	// NOTE this is only exposed for the purpose of testing
	KubeResourceStartTimeout time.Duration
	HelmJobStartTimeout      time.Duration
}

type Core struct {
	// Version is set by cmd/server/server.go
	Version string
	Settings

	// NOTE we set these 2 in cmd/server.go to prevent having to load all the
	// cloud provider various lib code everytime we load core
	AWSProvider  func(map[string]string) Provider
	DOProvider   func(map[string]string) Provider
	OSProvider   func(map[string]string) Provider
	GCEProvider  func(map[string]string) Provider
	PACKProvider func(map[string]string) Provider

	// A little different from the above
	K8SProvider Provider

	K8S        func(*model.Kube) kubernetes.ClientInterface
	KubeClient func(*model.Kube) (kclient.Interface, error)

	DefaultProvisioner Provisioner

	APIClient func(authType string, authToken string) *client.Client

	Log *logrus.Logger

	DB DBInterface

	Sessions      SessionsInterface
	Users         *Users
	CloudAccounts *CloudAccounts
	Kubes         *Kubes
	KubeResources KubeResourcesInterface
	Nodes         NodesInterface
	LoadBalancers *LoadBalancers
	HelmRepos     *HelmRepos
	HelmCharts    *HelmCharts
	HelmReleases  *HelmReleases

	// TODO should this be a pseudo-collection like Sessions?
	Actions *SafeMap
}

// NOTE this used to be core.New(), but due to how we load in values from the
// cli package, I needed to first actually initialize a Core struct and then
// configure.
func (c *Core) Initialize() {
	if err := c.InitializeForeground(); err != nil {
		panic(err)
	}
	if err := c.detectOrCreateAdmin(); err != nil {
		panic(err)
	}
	if err := c.detectOrCreateStableHelmRepo(); err != nil {
		panic(err)
	}
	c.InitializeBackground()
}

//------------------------------------------------------------------------------

// InitializeForeground sets up Log and DB on *Core.
func (c *Core) InitializeForeground() error {
	if c.ConfigFilePath != "" {
		var configFileSettings Settings
		configFile, err := os.Open(c.ConfigFilePath)
		if err != nil {
			return err
		}
		if err := json.NewDecoder(configFile).Decode(&configFileSettings); err != nil {
			return err
		}
		// Merge in command line settings (which overwrite respective config file settings)
		if err := mergo.Merge(&c.Settings, configFileSettings); err != nil {
			return err
		}
		// Set Default Settings with struct tags
		if err := defaults.Set(&c.Settings); err != nil {
			return err
		}
	}

	requiredFlags := map[string]string{
		"publish-host": c.PublishHost,
		"http-port":    c.HTTPPort,
	}
	for flag, val := range requiredFlags {
		if val == "" {
			return errors.New(flag + " required")
		}
	}

	// Log to stdout by default instead of stderr
	c.Log = logrus.New()
	c.Log.Out = os.Stdout

	if c.LogLevel != "" {
		levelInt, err := logrus.ParseLevel(c.LogLevel)
		if err != nil {
			return err
		}
		c.Log.SetLevel(levelInt)
		logrus.SetLevel(levelInt)
	}
	// db.LogMode(true)

	// DB
	var gormDB *gorm.DB
	var err error

	if c.PsqlHost != "" && c.PsqlDb != "" && c.PsqlUser != "" && c.PsqlPass != "" {
		// Postgres
		options := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", c.PsqlHost, c.PsqlDb, c.PsqlUser, c.PsqlPass)
		gormDB, err = gorm.Open("postgres", options)

	} else if c.SQLiteFile != "" {
		// SQLite3
		gormDB, err = gorm.Open("sqlite3", c.SQLiteFile)

	} else {
		err = errors.New("No DB configured. Must provide --psql-* options or --sqlite-file.")
	}
	if err != nil {
		return err
	}

	err = gormDB.AutoMigrate(
		&model.User{},
		&model.Kube{},
		&model.KubeResource{},
		&model.CloudAccount{},
		&model.Node{},
		&model.LoadBalancer{},
		&model.HelmRepo{},
		&model.HelmChart{},
		&model.HelmRelease{},
	).Error
	if err != nil {
		return err
	}

	c.DB = &DB{c, gormDB}

	c.Users = &Users{Collection{c}}
	c.Kubes = &Kubes{Collection{c}}
	c.KubeResources = &KubeResources{Collection{c}}
	c.CloudAccounts = &CloudAccounts{Collection{c}}
	c.Nodes = &Nodes{Collection{c}}
	c.LoadBalancers = &LoadBalancers{Collection{c}}
	c.HelmRepos = &HelmRepos{Collection{c}}
	c.HelmCharts = &HelmCharts{Collection{c}}
	c.HelmReleases = &HelmReleases{Collection{c}}
	c.Sessions = NewSessions(c)

	// Actions for async work
	c.Actions = NewSafeMap(c)

	// Kubernetes Client
	c.K8S = func(kube *model.Kube) kubernetes.ClientInterface {
		return kubernetes.NewClient(kube, kubernetes.DefaultHTTPClient)
	}
	c.KubeClient = func(kube *model.Kube) (kclient.Interface, error) {
		if kube.KubeAPIPort == "" {
			kube.KubeAPIPort = "443"
		}

		config, err := kubernetes.BuildBasicAuthConfig(kube.MasterPublicIP, kube.KubeAPIPort, kube.Username, kube.Password)
		if err != nil {
			return nil, errors.Wrap(err, "build config")
		}

		return kclient.NewForConfig(config)
	}

	// Kubernetes Provisioners
	c.DefaultProvisioner = &DefaultProvisioner{c}

	c.KubeResourceStartTimeout = 20 * time.Minute
	c.HelmJobStartTimeout = 30 * time.Second

	// API Client
	c.APIClient = func(authType string, authToken string) *client.Client {
		client := client.New(c.BaseURL(), authType, authToken, c.SSLCertFile)
		client.Version = c.Version
		return client
	}

	return nil
}

// InitializeBackground starts Action processing and RecurringServices for *Core.
func (c *Core) InitializeBackground() {
	// Recurring services
	if c.CapacityServiceEnabled {
		capacityService := &RecurringService{
			core: c,
			service: &CapacityService{
				Core:            c,
				WaitBeforeScale: 2 * time.Minute,
			},
			interval: 30 * time.Second,
			tag:      "Capacity Service",
		}
		go capacityService.Run()
	}

	nodeObserver := &RecurringService{
		core:     c,
		service:  &NodeObserver{c},
		interval: 30 * time.Second,
		tag:      "Node Observer",
	}
	go nodeObserver.Run()

	// TODO we probably don't need both this and the observer
	kubeResourcePopulater := &RecurringService{
		core:     c,
		fn:       c.KubeResources.Populate,
		interval: 30 * time.Second,
		tag:      "Kube Resource Populator",
	}
	go kubeResourcePopulater.Run()

	helmChartPopulater := &RecurringService{
		core:     c,
		fn:       c.HelmCharts.Populate,
		interval: 60 * time.Second,
		tag:      "Helm Chart Populator",
	}
	go helmChartPopulater.Run()

	helmReleasePopulater := &RecurringService{
		core:     c,
		fn:       c.HelmReleases.Populate,
		interval: 30 * time.Second,
		tag:      "Helm Release Populator",
	}
	go helmReleasePopulater.Run()

	sessionExpirer := &RecurringService{
		core:     c,
		service:  &SessionExpirer{c},
		interval: 15 * time.Second,
		tag:      "Session Expire",
	}
	go sessionExpirer.Run()
}

//------------------------------------------------------------------------------

func (c *Core) SSLEnabled() bool {
	return c.HTTPSPort != "" && c.SSLCertFile != "" && c.SSLKeyFile != ""
}

func (c *Core) HTTPSURL() string {
	return fmt.Sprintf("https://%s:%s", c.PublishHost, c.HTTPSPort)
}

func (c *Core) HTTPURL() string {
	return fmt.Sprintf("http://%s:%s", c.PublishHost, c.HTTPPort)
}

func (c *Core) BaseURL() string {
	if c.SSLEnabled() {
		return c.HTTPSURL()
	}
	return c.HTTPURL()
}

func (c *Core) APIURL() string {
	return fmt.Sprintf("%s/api/v0", c.BaseURL())
}

func (c *Core) UIURL() string {
	return fmt.Sprintf("%s/ui", c.BaseURL())
}

//------------------------------------------------------------------------------

func (c *Core) detectOrCreateAdmin() error {
	if c.SupportPassword != "" {
		support := &model.User{
			Username: "support",
			Password: c.SupportPassword,
			Role:     model.UserRoleAdmin,
		}

		if err := c.DB.First(new(model.User), "username = ?", support.Username); err == nil {
			// Already have an support
			fmt.Println("Support User Already Configured...")
			return nil
		}
		if err := c.Users.Create(support); err != nil {
			return err
		}
	}
	if err := c.DB.First(new(model.User), "role = ?", model.UserRoleAdmin); err == nil {
		// Already have an admin
		return nil
	}

	// TODO this should be logged to file, but isn't due to how we open file in server.go
	c.Log.Info("No Admin detected, creating new and printing credentials:")

	password := util.RandomString(16)

	admin := &model.User{
		Username: "admin",
		Password: password,
		Role:     model.UserRoleAdmin,
	}
	if err := c.Users.Create(admin); err != nil {
		return err
	}

	// Print to STDOUT (not c.Log, which would save to file)
	fmt.Printf("\n  ( ͡° ͜ʖ ͡°)  USERNAME: admin  PASSWORD: %s\n\n", password)

	return nil
}

//------------------------------------------------------------------------------

func (c *Core) detectOrCreateStableHelmRepo() error {
	if err := c.DB.First(new(model.HelmRepo), "name = ?", "stable"); err == nil {
		return nil // Already exists
	}

	c.Log.Info("Registering stable HelmRepo")

	repo := &model.HelmRepo{
		Name: "stable",
		URL:  "https://kubernetes-charts.storage.googleapis.com",
	}
	return c.HelmRepos.Create(repo)
}
