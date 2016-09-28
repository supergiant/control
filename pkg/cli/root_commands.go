package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func Configure(c *cli.Context) error {
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

func Kubectl(c *cli.Context) error {
	id := c.Int64("kube-id")
	kube := new(model.Kube)
	if err := newClient(c).Kubes.Get(&id, kube); err != nil {
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
