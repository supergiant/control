package cli

import (
	"os"

	"github.com/mitchellh/go-homedir"
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

func Run() {
	app := cli.NewApp()

	app.Name = "supergiant"
	app.Usage = "Supergiant CLI"
	app.Version = "0.10.0-alpha" // TODO this shouldn't be hard-coded here

	app.Commands = []cli.Command{
		{
			Name:   "configure",
			Usage:  "globally configure server settings (helpful to prevent repeating flags)",
			Flags:  baseFlags,
			Action: Configure,
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
			Action: Kubectl,
		},
		{
			Name:  "cloud_accounts",
			Usage: "actions for Cloud Accounts",
			Subcommands: []cli.Command{
				commandList("Cloud Account", ListCloudAccounts),
				commandCreate("Cloud Account", CreateCloudAccount),
				commandGet("Cloud Account", GetCloudAccount),
				// commandUpdate("Cloud Account", UpdateCloudAccount), -- TODO ------ we need to define IMMUTABLE json tag to prevent update of fields like Name
				commandAction("delete", "Cloud Account", DeleteCloudAccount),
			},
		},
		{
			Name:  "entrypoints",
			Usage: "actions for Entrypoints",
			Subcommands: []cli.Command{
				commandList("Entrypoint", ListEntrypoints),
				commandCreate("Entrypoint", CreateEntrypoint),
				commandGet("Entrypoint", GetEntrypoint),
				// commandUpdate("Entrypoint", UpdateEntrypoint),
				commandAction("delete", "Entrypoint", DeleteEntrypoint),
			},
		},
		{
			Name:  "entrypoint_listeners",
			Usage: "actions for Entrypoint Listeners",
			Subcommands: []cli.Command{
				commandList("Entrypoint Listener", ListEntrypointListeners),
				commandCreate("Entrypoint Listener", CreateEntrypointListener),
				commandGet("Entrypoint Listener", GetEntrypointListener),
				// commandUpdate("Entrypoint Listener", UpdateEntrypointListener),
				commandAction("delete", "Entrypoint Listener", DeleteEntrypointListener),
			},
		},
		{
			Name:  "kubes",
			Usage: "actions for Kubes",
			Subcommands: []cli.Command{
				commandList("Kube", ListKubes),
				commandCreate("Kube", CreateKube),
				commandGet("Kube", GetKube),
				// commandUpdate("Kube", UpdateKube),
				commandAction("delete", "Kube", DeleteKube),
			},
		},
		{
			Name:  "nodes",
			Usage: "actions for Nodes",
			Subcommands: []cli.Command{
				commandList("Node", ListNodes),
				commandCreate("Node", CreateNode),
				commandGet("Node", GetNode),
				// commandUpdate("Node", UpdateNode),
				commandAction("delete", "Node", DeleteNode),
			},
		},
		{
			Name:  "sessions",
			Usage: "actions for Sessions",
			Subcommands: []cli.Command{
				commandList("Session", ListSessions),
				commandCreate("Session", CreateSession),
				commandGet("Session", GetSession),
				// commandUpdate("Session", UpdateSession),
				commandAction("delete", "Session", DeleteSession),
			},
		},
		{
			Name:  "users",
			Usage: "actions for Users",
			Subcommands: []cli.Command{
				commandList("User", ListUsers),
				commandCreate("User", CreateUser),
				commandGet("User", GetUser),
				commandUpdate("User", UpdateUser),
				commandAction("delete", "User", DeleteUser),
			},
		},
		{
			Name:  "volumes",
			Usage: "actions for Volumes",
			Subcommands: []cli.Command{
				commandList("Volume", ListVolumes),
				commandCreate("Volume", CreateVolume),
				commandGet("Volume", GetVolume),
				commandUpdate("Volume", UpdateVolume),
				commandAction("delete", "Volume", DeleteVolume),
			},
		},
		{
			Name:  "kube_resources",
			Usage: "actions for Kube Resources",
			Subcommands: []cli.Command{
				commandList("KubeResource", ListKubeResources),
				commandCreate("KubeResource", CreateKubeResource),
				commandGet("KubeResource", GetKubeResource),
				commandUpdate("KubeResource", UpdateKubeResource),
				commandAction("delete", "KubeResource", DeleteKubeResource),
				commandAction("start", "KubeResource", StartKubeResource),
				commandAction("stop", "KubeResource", StopKubeResource),
			},
		},
	}

	app.Run(os.Args)
}
