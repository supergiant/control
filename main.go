package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/supergiant/supergiant/api"
	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
)

func main() {

	app := cli.NewApp()
	app.Name = "supergiant-api"
	app.Usage = "The Supergiant api server."
	app.Action = func(ctx *cli.Context) {

		// Check the args. The ones we don't have default values for...
		if core.K8sUser == "<Kubernetes api userID>" {
			core.Log.Error("Kubernetes userID required...")
			cli.ShowCommandHelp(ctx, "")
			os.Exit(5)
		}
		if core.K8sPass == "<Kubernetes api password>" {
			core.Log.Error("Kubernetes Password required...")
			cli.ShowCommandHelp(ctx, "")
			os.Exit(5)
		}
		if core.AwsRegion == "" {
			region, err := getAWSRegion()
			if err != nil {
				log.Println("ERROR:", err)
			} else {
				core.AwsRegion = region
				log.Println("INFO: AWS Region Detected,", core.AwsRegion)
			}
		}
		if core.AwsAZ == "" {
			az, err := getAWSAZ()
			if err != nil {
				log.Println("ERROR:", err)
			} else {
				core.AwsAZ = az
				log.Println("INFO: AWS AZ Detected,", core.AwsAZ)
			}
		}
		if core.AwsSgID == "" {
			sg, err := getAWSSecurityGroupID()
			if err != nil {
				log.Println("ERROR:", err)
			} else {
				core.AwsSgID = sg
				log.Println("INFO: AWS Security Group Detected,", core.AwsSgID)
			}
		}
		if core.AwsSubnetID == "" {
			sub, err := getAWSSubnetID()
			if err != nil {
				log.Println("ERROR:", err)
			} else {
				core.AwsSubnetID = sub
				log.Println("INFO: AWS Security Group Detected,", core.AwsSubnetID)
			}
		}

		core.EtcdEndpoints = ctx.StringSlice("etcd-host")
		if len(core.EtcdEndpoints) < 0 {
			core.EtcdEndpoints = []string{"http://etcd:2379"}
		}

		// Log args that have default values.
		core.Log.Info("ETCD hosts,", ctx.StringSlice("etcd-host"))
		core.Log.Info("Kubernetes Host,", core.K8sHost)

		c := core.New(
			ctx.Bool("https-mode"),   // Tells the api if it needs to connect to Kuberntes over TLS or not.
			ctx.String("access-key"), // AWS Access Key
			ctx.String("secret-key"), // AWS Secret Key
		)

		// TODO should probably be able to say api.New(), because we shouldn't have to import task here
		// NOTE using pool size of 4
		go task.NewSupervisor(c, 20).Run()

		router := api.NewRouter(c)

		core.Log.Info("Serving API on port :8080")
		core.Log.Info(http.ListenAndServe(":8080", router))
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "etcd-host",
			Usage:  "Array of etcd hosts.",
			EnvVar: "ETCD_ENDPOINT",
		},
		cli.StringFlag{
			Name:        "k8sHost, kh",
			Value:       "kubernetes",
			Usage:       "IP of a Kuberntes api.",
			EnvVar:      "K8S_HOST",
			Destination: &core.K8sHost,
		},
		cli.StringFlag{
			Name:        "k8sUser, ku",
			Value:       "<Kubernetes api userID>",
			Usage:       "Username used to connect to your Kubernetes api.",
			EnvVar:      "K8S_USER",
			Destination: &core.K8sUser,
		},
		cli.StringFlag{
			Name:        "k8sPass, kp",
			Value:       "<Kubernetes api password>",
			Usage:       "Password used to connect to your Kubernetes api.",
			EnvVar:      "K8S_PASS",
			Destination: &core.K8sPass,
		},
		cli.StringFlag{
			Name:        "awsRegion, ar",
			Value:       "",
			Usage:       "AWS Region in which your kubernetes cluster resides.",
			EnvVar:      "AWS_REGION",
			Destination: &core.AwsRegion,
		},
		cli.StringFlag{
			Name:        "awsAZ, az",
			Value:       "",
			Usage:       "AWS Availability Zone in which your kubernetes cluster resides.",
			EnvVar:      "AWS_AZ",
			Destination: &core.AwsAZ,
		},
		cli.StringFlag{
			Name:        "awsSgID, sg",
			Value:       "",
			Usage:       "AWS Security Group in which your kubernetes cluster resides.",
			EnvVar:      "AWS_SG_ID",
			Destination: &core.AwsSgID,
		},
		cli.StringFlag{
			Name:        "awsSubnetID, sid",
			Value:       "",
			Usage:       "AWS Subnet ID in which your kubernetes cluster resides.",
			EnvVar:      "AWS_SUBNET_ID",
			Destination: &core.AwsSubnetID,
		},
		cli.StringFlag{
			Name:  "access-key",
			Value: "",
			Usage: "AWS Access key.",
		},
		cli.StringFlag{
			Name:  "secret-key",
			Value: "",
			Usage: "AWS Secret key.",
		},
		cli.BoolFlag{
			Name:   "https-mode",
			Usage:  "Enable https mode for guber client when running outside a kube.",
			EnvVar: "HTTPS_MODE",
		},
	}

	app.Run(os.Args)
}
