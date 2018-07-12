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
	port         = flag.Int("port", 8080, "tcp port to listen for incoming requests")
	etcdURL      = flag.String("etcd_url", "localhost:2379", "etcd url with port")
	templatesDir = flag.String("templates", "/etc/supergiant/templates", "supergiant will load script templates from the specified directory on start")
	logLevel     = flag.String("log_level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
)

func main() {
	flag.Parse()

	cfg := &controlplane.Config{
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
		_ = <-sigs
		logrus.Info("shutting down...")
		server.Stop()
	}()

	logrus.Infof("supergiant is starting on port %d", *port)
	server.Start()
}
