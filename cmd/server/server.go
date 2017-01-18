package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/pkg/provider/digitalocean"
	"github.com/supergiant/supergiant/pkg/provider/gce"
	"github.com/supergiant/supergiant/pkg/provider/openstack"
	"github.com/supergiant/supergiant/pkg/server"
)

var version = "unversioned"

func main() {
	c := &core.Core{
		Version: version,
	}

	app := cli.NewApp()
	app.Name = "supergiant-server"
	app.Usage = "Supergiant server " + version

	app.Action = func(ctx *cli.Context) {

		// TODO should check for missing setting here to show cli help
		c.Initialize()

		// See relevant NOTE in core.go
		c.AWSProvider = func(creds map[string]string) core.Provider {
			return &aws.Provider{
				Core: c,
				EC2:  aws.EC2,
				IAM:  aws.IAM,
				ELB:  aws.ELB,
				S3:   aws.S3,
			}
		}
		c.DOProvider = func(creds map[string]string) core.Provider {
			return &digitalocean.Provider{
				Core:   c,
				Client: digitalocean.Client,
			}
		}
		c.OSProvider = func(creds map[string]string) core.Provider {
			return &openstack.Provider{
				Core:   c,
				Client: openstack.Client,
			}
		}
		c.GCEProvider = func(creds map[string]string) core.Provider {
			return &gce.Provider{
				Core:   c,
				Client: gce.Client,
			}
		}

		// We do this here, and not in core, so that we can ensure the file closes on exit.
		if c.LogPath != "" {
			file, err := os.OpenFile(c.LogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			path, err := filepath.Abs(c.LogPath)
			if err != nil {
				panic(err)
			}
			fmt.Println("Writing log to " + path)
			c.Log.Out = file
		}

		srv, err := server.New(c)
		if err != nil {
			panic(err)
		}

		srv.Start()
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "sqlite-file",
			Usage:       "SQLite3 database (.db) file",
			Destination: &c.SQLiteFile,
		},
		cli.StringFlag{
			Name:        "psql-host",
			Usage:       "PostgreSQL host",
			Destination: &c.PsqlHost,
		},
		cli.StringFlag{
			Name:        "psql-db",
			Usage:       "PostgreSQL database name",
			Destination: &c.PsqlDb,
		},
		cli.StringFlag{
			Name:        "psql-user",
			Usage:       "PostgreSQL database user",
			Destination: &c.PsqlUser,
		},
		cli.StringFlag{
			Name:        "psql-pass",
			Usage:       "PostgreSQL database password",
			Destination: &c.PsqlPass,
		},
		cli.StringFlag{
			Name:        "publish-host",
			Usage:       "Host that can be used to connect to this Supergiant server remotely",
			Destination: &c.PublishHost,
		},
		cli.StringFlag{
			Name:        "http-port",
			Usage:       "HTTP port for the web interfaces",
			Destination: &c.HTTPPort,
		},
		cli.BoolFlag{
			Name:        "ui-enabled",
			Usage:       "Enabled UI",
			Destination: &c.UIEnabled,
		},
		cli.StringFlag{
			Name:        "https-port",
			Usage:       "HTTPS (SSL) port for the web interfaces",
			Destination: &c.HTTPSPort,
		},
		cli.StringFlag{
			Name:        "ssl-cert-file",
			Usage:       "SSL certificate file",
			Destination: &c.SSLCertFile,
		},
		cli.StringFlag{
			Name:        "ssl-key-file",
			Usage:       "SSL key file",
			Destination: &c.SSLKeyFile,
		},
		cli.StringFlag{
			Name:        "log-file",
			Usage:       "Log output filepath",
			Destination: &c.LogPath,
		},
		cli.StringFlag{
			Name:        "log-level",
			Usage:       "Log level",
			Destination: &c.LogLevel,
			// Value:  <--- NOTE just cuz you always forget you can set defaults
		},
		cli.StringFlag{
			Name:        "config-file",
			Usage:       "JSON config filepath (command line arguments will override the values set here)",
			Destination: &c.ConfigFilePath,
		},
	}

	app.Run(os.Args)
}
