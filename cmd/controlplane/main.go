package main

import (
	"flag"

	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/controlplane"
)

var (
	addr         = flag.String("address", "0.0.0.0", "network interface to attach server to")
	port         = flag.Int("port", 8080, "tcp port to listen for incoming requests")
	etcdURL      = flag.String("etcd-url", "localhost:2379", "etcd url with port")
	templatesDir = flag.String("templates", "/etc/supergiant/templates", "supergiant will load script templates from the specified directory on start")
	logLevel     = flag.String("log-level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
)

func main() {
	flag.Parse()

	cfg := &controlplane.Config{
		Addr:         *addr,
		Port:         *port,
		EtcdUrl:      *etcdURL,
		TemplatesDir: *templatesDir,
		LogLevel:     *logLevel,
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
