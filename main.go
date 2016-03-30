package main

import (
	"fmt"
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
	app.Action = func(c *cli.Context) {

		// Check the args. The ones we don't have default values for...
		if core.K8sUser == "<Kubernetes api userID>" {
			fmt.Println("ERROR: Kubernetes userID required...")
			cli.ShowCommandHelp(c, "")
			os.Exit(5)
		}
		if core.K8sPass == "<Kubernetes api password>" {
			fmt.Println("ERROR: Kubernetes Password required...")
			cli.ShowCommandHelp(c, "")
			os.Exit(5)
		}
		if core.AwsRegion == "<AWS Region>" {
			fmt.Println("ERROR: AWS Region required...")
			cli.ShowCommandHelp(c, "")
			os.Exit(5)
		}
		if core.AwsAZ == "<AWS Availability Zone>" {
			fmt.Println("ERROR: AWS Availability Zone required...")
			cli.ShowCommandHelp(c, "")
			os.Exit(5)
		}
		core.EtcdEndpoints = c.StringSlice("etcd-host")
		if len(core.EtcdEndpoints) < 0 {
			core.EtcdEndpoints = []string{"http://etcd:2379"}
		}

		// Log args that have default values.
		log.Println("INFO: ETCD hosts,", c.StringSlice("etcd-host"))
		log.Println("INFO: Kubernetes Host,", core.K8sHost)

		core := core.New()

		// TODO should probably be able to say api.New(), because we shouldn't have to import task here
		// NOTE using pool size of 4
		go task.NewSupervisor(core, 20).Run()

		router := api.NewRouter(core)

		log.Println("INFO: Serving API on port :8080")
		log.Fatal(http.ListenAndServe(":8080", router))
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
			Value:       "<AWS Region>",
			Usage:       "AWS Region in which your kubernetes cluster resides.",
			EnvVar:      "AWS_REGION",
			Destination: &core.AwsRegion,
		},
		cli.StringFlag{
			Name:        "awsAZ, az",
			Value:       "<AWS Availability Zone>",
			Usage:       "AWS Availability Zone in which your kubernetes cluster resides.",
			EnvVar:      "AWS_AZ",
			Destination: &core.AwsAZ,
		},
		cli.StringFlag{
			Name:        "awsSgID, sg",
			Value:       "<AWS Security Group ID>",
			Usage:       "AWS Security Group in which your kubernetes cluster resides.",
			EnvVar:      "AWS_SG_ID",
			Destination: &core.AwsSgID,
		},
		cli.StringFlag{
			Name:        "awsSubnetID, sid",
			Value:       "<AWS Subnet ID>",
			Usage:       "AWS Subnet ID in which your kubernetes cluster resides.",
			EnvVar:      "AWS_SUBNET_ID",
			Destination: &core.AwsSubnetID,
		},
	}

	app.Run(os.Args)
}
