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
	version       = "unstable"
	addr          = flag.String("address", "0.0.0.0", "network interface to attach server to")
	port          = flag.Int("port", 0, "secure tcp port to listen to for incoming HTTPS requests. Provide server certificates with -cert-file and -key-file flags")
	insecurePort  = flag.Int("insecure-port", 8080, "tcp port to listen for incoming HTTP requests. if -port is set this flag will be ignored")
	certFile      = flag.String("cert-file", "", "file containing server x509 certificate")
	keyFile       = flag.String("key-file", "", "file containing x509 private key matching --cert-file")
	storageMode   = flag.String("storage-mode", "file", "storage type either file(default), memory or etcd")
	storageURI    = flag.String("storage-uri", "supergiant.db", "uri of storage depends on selected storage type, for memory storage type this is empty")
	templatesDir  = flag.String("templates", "", "supergiant will load script templates from the specified directory on start")
	logDir        = flag.String("log-dir", "/tmp", "logging directory for task logs")
	logLevel      = flag.String("log-level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
	logFormat     = flag.String("log-format", "txt", "logging format [txt json]")
	spawnInterval = flag.Int("spawnInterval", 5, "interval between API calls to cloud provider for creating instance")
	//TODO: rewrite to single flag port-range
	ProxiesPortRangeFrom = flag.Int("proxies-port-from", 60200, "first tcp port in a range of binding reverse proxies for service apps")
	ProxiesPortRangeTo   = flag.Int("proxies-port-to", 60250, "last tcp port in a range of binding reverse proxies for service apps")
	pprofListenStr       = flag.String("pprofListenStr", "",
		"pprof listen str host:port")
)

func main() {
	flag.Parse()

	configureLogging(*logLevel, *logFormat)

	cfg := &controlplane.Config{
		Addr:          *addr,
		Port:          *port,
		InsecurePort:  *insecurePort,
		CertFile:      *certFile,
		KeyFile:       *keyFile,
		StorageMode:   *storageMode,
		StorageURI:    *storageURI,
		TemplatesDir:  *templatesDir,
		LogDir:        *logDir,
		ReadTimeout:   time.Second * 60,
		WriteTimeout:  time.Second * 30,
		IdleTimeout:   time.Second * 120,
		SpawnInterval: time.Second * time.Duration(*spawnInterval),

		PprofListenStr: *pprofListenStr,

		ProxiesPortRange: proxy.PortRange{int32(*ProxiesPortRangeFrom), int32(*ProxiesPortRangeTo)},
		Version:          version,
	}

	server, err := controlplane.New(cfg)
	if err != nil {
		logrus.Infof("configuration: %+v", *cfg)
		logrus.Fatalf("broken configuration: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logrus.Info("shutting down...")
		server.Shutdown()
	}()

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
