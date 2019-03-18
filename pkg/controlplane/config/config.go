package config

import (
	"time"
	"github.com/supergiant/control/pkg/proxy"
)

// Config is the server configuration
type Config struct {
	Port          int
	Addr          string
	StorageMode   string
	StorageURI    string
	TemplatesDir  string
	UIDir         string
	SpawnInterval time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	PprofListenStr string

	ProxiesPortRange proxy.PortRange

	Version string
}

