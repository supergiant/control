package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/apiserver"
	"os"
	"os/signal"
	"syscall"
)

var (
	port         = flag.Int("port", 8080, "tcp port to listen for incoming requests")
	etcdURL      = flag.String("etcdURL", "localhost:2379", "etcd url with port")
	templatesDir = flag.String("scirpts_dir", "/etc/supergiant/scripts", "supergiant will load script templates from the specified directory on start")
	logLevel     = flag.String("log_level", "INFO", "logging level, e.g. info, warning, debug, error, fatal")
)

func main() {
	flag.Parse()

	cfg := &apiserver.Config{
		Port:         *port,
		EtcdUrl:      *etcdURL,
		TemplatesDir: *templatesDir,
		LogLevel:     *logLevel,
	}

	s, err := apiserver.New(cfg)
	if err != nil {
		logrus.Fatalf("broken configuration: %v", err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		logrus.Info("shutting down...")
		s.Stop()
	}()
	s.Start()
}
