package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/controlplane"
	"github.com/supergiant/control/pkg/proxy"
)

var (
	version       = "undefined"
	addr          = flag.String("address", "0.0.0.0", "network interface to attach server to")
	port          = flag.Int("port", 8080, "tcp port to listen for incoming requests")
	etcdURL       = flag.String("etcd-url", "localhost:2379", "etcd url with port")
	templatesDir  = flag.String("templates", "/etc/supergiant/templates/", "supergiant will load script templates from the specified directory on start")
	logLevel      = flag.String("log-level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
	logFormat     = flag.String("log-format", "txt", "logging format [txt json]")
	spawnInterval = flag.Int("spawnInterval", 5, "interval between API calls to cloud provider for creating instance")
	uiDir         = flag.String("ui-dir", "./cmd/ui/assets/dist", "directory for supergiant ui to be served")
	//TODO: rewrite to single flag port-range
	ProxiesPortRangeFrom = flag.Int("proxies-port-from", 60200, "first tcp port in a range of binding reverse proxies for service apps")
	ProxiesPortRangeTo   = flag.Int("proxies-port-to", 60250, "last tcp port in a range of binding reverse proxies for service apps")
	pprofListenStr = flag.String("pprofListenStr", "",
		"pprof listen str host:port")
)

func main() {
	flag.Parse()

	configureLogging(*logLevel, *logFormat)

	cfg := &controlplane.Config{
		Addr:             *addr,
		Port:             *port,
		EtcdUrl:          *etcdURL,
		TemplatesDir:     *templatesDir,
		ReadTimeout:      time.Second * 20,
		WriteTimeout:     time.Second * 10,
		IdleTimeout:      time.Second * 120,
		SpawnInterval:    time.Second * time.Duration(*spawnInterval),
		UiDir:            *uiDir,

		PprofListenStr: *pprofListenStr,

		ProxiesPortRange: proxy.PortRange{int32(*ProxiesPortRangeFrom), int32(*ProxiesPortRangeTo)},
		Version: version,
	}

	server, err := controlplane.New(cfg)
	if err != nil {
		logrus.Fatalf("broken configuration: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logrus.Info("shutting down...")
		server.Shutdown()
	}()

	logrus.Infof("supergiant is starting on port %d", *port)
	server.Start()
}

// TODO: create sglog package
func configureLogging(level, format string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		// set logLevel to INFO by default
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)

	switch strings.TrimSpace(format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}
