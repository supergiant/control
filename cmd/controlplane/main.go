package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/controlplane"
)

var (
	addr         = flag.String("address", "0.0.0.0", "network interface to attach server to")
	port         = flag.Int("port", 8080, "tcp port to listen for incoming requests")
	etcdURL      = flag.String("etcd-url", "localhost:2379", "etcd url with port")
	templatesDir = flag.String("templates", "/etc/supergiant/templates/", "supergiant will load script templates from the specified directory on start")
	logLevel     = flag.String("log-level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
	logFormat    = flag.String("log-format", "txt", "logging format [txt json]")
)

func main() {
	flag.Parse()

	configureLogging(*logLevel, *logFormat)

	cfg := &controlplane.Config{
		Addr:         *addr,
		Port:         *port,
		EtcdUrl:      *etcdURL,
		TemplatesDir: *templatesDir,
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
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
