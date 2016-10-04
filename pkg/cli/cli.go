package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"

	"github.com/mitchellh/go-homedir"
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

var globalConfFile string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	globalConfFile = home + "/.supergiant"
}

var baseFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "server, s",
		Usage: "Host and port of the Supergiant server",
	},
	cli.StringFlag{
		Name:  "api-token, t",
		Usage: "API token of the operating Supergiant User",
	},
	cli.StringFlag{
		Name:  "cert-file, c",
		Usage: "Filepath of the SSL certificate used by the server. If not provided, the cert must be manually trusted through OS.",
	},
}

type CLI struct {
	*cli.App
	Client  func(*cli.Context) *client.Client
	Stdin   *os.File
	Version string
}

func New(clientFn func(*cli.Context) *client.Client, stdin *os.File, version string) *CLI {
	sgcli := &CLI{cli.NewApp(), clientFn, stdin, version}

	sgcli.Name = "supergiant"
	sgcli.Usage = "Supergiant CLI " + version
	// TODO for whatever reason the version always reports 0.0.0
	sgcli.Version = version

	sgcli.Commands = []cli.Command{
		{
			Name:   "configure",
			Usage:  "globally configure server settings (helpful to prevent repeating flags)",
			Flags:  baseFlags,
			Action: sgcli.commandConfigure,
		},
		{
			Name:  "kubectl",
			Usage: "wrapper for Kubectl that auto-populates connection-related tags",
			Flags: append(baseFlags, []cli.Flag{
				cli.Int64Flag{
					Name:  "kube-id, k",
					Usage: "ID of the Supergiant Kube",
				},
			}...),
			Action: sgcli.commandKubectl,
		},
		{
			Name:  "cloud_accounts",
			Usage: "actions for CloudAccounts",
			Subcommands: []cli.Command{
				sgcli.commandList("CloudAccounts", new(model.CloudAccountList)),
				sgcli.commandCreate("CloudAccounts", new(model.CloudAccount)),
				sgcli.commandGet("CloudAccounts", new(model.CloudAccount)),
				sgcli.commandUpdate("CloudAccounts", new(model.CloudAccount)),
				sgcli.commandAction("delete", "Delete", "CloudAccounts", new(model.CloudAccount)),
			},
		},
		{
			Name:  "entrypoints",
			Usage: "actions for Entrypoints",
			Subcommands: []cli.Command{
				sgcli.commandList("Entrypoints", new(model.EntrypointList)),
				sgcli.commandCreate("Entrypoints", new(model.Entrypoint)),
				sgcli.commandGet("Entrypoints", new(model.Entrypoint)),
				sgcli.commandUpdate("Entrypoints", new(model.Entrypoint)),
				sgcli.commandAction("delete", "Delete", "Entrypoints", new(model.Entrypoint)),
			},
		},
		{
			Name:  "entrypoint_listeners",
			Usage: "actions for EntrypointListeners",
			Subcommands: []cli.Command{
				sgcli.commandList("EntrypointListeners", new(model.EntrypointListenerList)),
				sgcli.commandCreate("EntrypointListeners", new(model.EntrypointListener)),
				sgcli.commandGet("EntrypointListeners", new(model.EntrypointListener)),
				sgcli.commandUpdate("EntrypointListeners", new(model.EntrypointListener)),
				sgcli.commandAction("delete", "Delete", "EntrypointListeners", new(model.EntrypointListener)),
			},
		},
		{
			Name:  "kubes",
			Usage: "actions for Kubes",
			Subcommands: []cli.Command{
				sgcli.commandList("Kubes", new(model.KubeList)),
				sgcli.commandCreate("Kubes", new(model.Kube)),
				sgcli.commandGet("Kubes", new(model.Kube)),
				sgcli.commandUpdate("Kubes", new(model.Kube)),
				sgcli.commandAction("delete", "Delete", "Kubes", new(model.Kube)),
			},
		},
		{
			Name:  "nodes",
			Usage: "actions for Nodes",
			Subcommands: []cli.Command{
				sgcli.commandList("Nodes", new(model.NodeList)),
				sgcli.commandCreate("Nodes", new(model.Node)),
				sgcli.commandGet("Nodes", new(model.Node)),
				sgcli.commandUpdate("Nodes", new(model.Node)),
				sgcli.commandAction("delete", "Delete", "Nodes", new(model.Node)),
			},
		},
		{
			Name:  "sessions",
			Usage: "actions for Sessions",
			Subcommands: []cli.Command{
				sgcli.commandList("Sessions", new(model.SessionList)),
				sgcli.commandCreate("Sessions", new(model.Session)),
				sgcli.commandGet("Sessions", new(model.Session)),
				sgcli.commandAction("delete", "Delete", "Sessions", new(model.Session)),
			},
		},
		{
			Name:  "users",
			Usage: "actions for Users",
			Subcommands: []cli.Command{
				sgcli.commandList("Users", new(model.UserList)),
				sgcli.commandCreate("Users", new(model.User)),
				sgcli.commandGet("Users", new(model.User)),
				sgcli.commandUpdate("Users", new(model.User)),
				sgcli.commandAction("delete", "Delete", "Users", new(model.User)),
			},
		},
		{
			Name:  "volumes",
			Usage: "actions for Volumes",
			Subcommands: []cli.Command{
				sgcli.commandList("Volumes", new(model.VolumeList)),
				sgcli.commandCreate("Volumes", new(model.Volume)),
				sgcli.commandGet("Volumes", new(model.Volume)),
				sgcli.commandUpdate("Volumes", new(model.Volume)),
				sgcli.commandAction("delete", "Delete", "Volumes", new(model.Volume)),
				sgcli.commandAction("resize", "Resize", "Volumes", new(model.Volume)),
			},
		},
		{
			Name:  "kube_resources",
			Usage: "actions for Kube Resources",
			Subcommands: []cli.Command{
				sgcli.commandList("KubeResources", new(model.KubeResourceList)),
				sgcli.commandCreate("KubeResources", new(model.KubeResource)),
				sgcli.commandGet("KubeResources", new(model.KubeResource)),
				sgcli.commandUpdate("KubeResources", new(model.KubeResource)),
				sgcli.commandAction("delete", "Delete", "KubeResources", new(model.KubeResource)),
				sgcli.commandAction("start", "Start", "KubeResources", new(model.KubeResource)),
				sgcli.commandAction("stop", "Stop", "KubeResources", new(model.KubeResource)),
			},
		},
	}

	return sgcli
}

// Private

func (sgcli *CLI) commandList(collectionName string, list model.List) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list " + collectionName,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: "--filter=name:this,or,that --filter=other_field:value",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "--format=\"{{ .ThisField }}\"",
			},
		}...),
		Action: func(c *cli.Context) error {
			filters, err := listFilters(c)
			if err != nil {
				return err
			}
			list.SetFilters(filters)

			fn := reflect.ValueOf(sgcli.Client(c)).Elem().FieldByName(collectionName).MethodByName("List")
			ret := fn.Call([]reflect.Value{reflect.ValueOf(list)})

			if err := ret[0].Interface(); err != nil {
				return err.(error)
			}
			return printList(c, list)
		},
	}
}

func (sgcli *CLI) commandCreate(collectionName string, item model.Model) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create new " + collectionName,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "JSON input file",
			},
		}...),
		Action: func(c *cli.Context) error {
			if err := sgcli.decodeInputFileInto(c, item); err != nil {
				return err
			}

			fn := reflect.ValueOf(sgcli.Client(c)).Elem().FieldByName(collectionName).MethodByName("Create")
			ret := fn.Call([]reflect.Value{reflect.ValueOf(item)})
			if err := ret[0].Interface(); err != nil {
				return err.(error)
			}

			return printObj(item)
		},
	}
}

func (sgcli *CLI) commandGet(collectionName string, item model.Model) cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get " + collectionName,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
		}...),
		Action: func(c *cli.Context) error {
			id := c.Int64("id")

			fn := reflect.ValueOf(sgcli.Client(c)).Elem().FieldByName(collectionName).MethodByName("Get")
			ret := fn.Call([]reflect.Value{reflect.ValueOf(&id), reflect.ValueOf(item)})
			if err := ret[0].Interface(); err != nil {
				return err.(error)
			}

			return printObj(item)
		},
	}
}

func (sgcli *CLI) commandUpdate(collectionName string, item model.Model) cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update " + collectionName,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
			cli.StringFlag{
				Name:  "file, f",
				Usage: "JSON input file",
			},
		}...),
		Action: func(c *cli.Context) error {
			id := c.Int64("id")
			if err := sgcli.decodeInputFileInto(c, item); err != nil {
				return err
			}

			fn := reflect.ValueOf(sgcli.Client(c)).Elem().FieldByName(collectionName).MethodByName("Update")
			ret := fn.Call([]reflect.Value{reflect.ValueOf(&id), reflect.ValueOf(item)})
			if err := ret[0].Interface(); err != nil {
				return err.(error)
			}

			return printObj(item)
		},
	}
}

func (sgcli *CLI) commandAction(action string, methodName string, collectionName string, item model.Model) cli.Command {
	return cli.Command{
		Name:  action,
		Usage: action + " " + collectionName,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
		}...),
		Action: func(c *cli.Context) error {
			id := c.Int64("id")

			fn := reflect.ValueOf(sgcli.Client(c)).Elem().FieldByName(collectionName).MethodByName(methodName)
			ret := fn.Call([]reflect.Value{reflect.ValueOf(&id), reflect.ValueOf(item)})
			if err := ret[0].Interface(); err != nil {
				return err.(error)
			}

			return printObj(item)
		},
	}
}

// Root commands

func (sgcli *CLI) commandConfigure(c *cli.Context) error {
	conf := newGlobalConf(c)
	b, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(globalConfFile, b, 0600); err != nil {
		return err
	}
	fmt.Println("Written to " + globalConfFile)
	return nil
}

func (sgcli *CLI) commandKubectl(c *cli.Context) error {
	id := c.Int64("kube-id")
	kube := new(model.Kube)
	if err := sgcli.Client(c).Kubes.Get(&id, kube); err != nil {
		return err
	}

	args := []string{
		"--api-version=v1",
		"--insecure-skip-tls-verify=true",
		"--server=https://" + kube.MasterPublicIP,
		"--cluster=" + kube.Name,
		"--username=" + kube.Username,
		"--password=" + kube.Password,
	}
	// Prepend the user's kubectl args, which will look like `get pods --namespace=kube-system`
	args = append([]string(c.Args()), args...)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Helpers

func (sgcli *CLI) decodeInputFileInto(c *cli.Context, item model.Model) (err error) {
	var file *os.File

	switch filepath := c.String("f"); filepath {
	case "":
		return errors.New("-f required")
	case "-":
		file = sgcli.Stdin
	default:
		file, err = os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	err = json.NewDecoder(file).Decode(item)
	return err
}
