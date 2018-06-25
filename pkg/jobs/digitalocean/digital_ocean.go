package digitalocean

import (
	"strings"
	"github.com/supergiant/supergiant/pkg/util"
	"bytes"
	"github.com/digitalocean/godo"
	"strconv"
	"time"
)

type Job struct{

}

func (j *Job) CreateDroplet(name, k8sVersion, region, size string, fingerprints []string) error {
	name = name + "-master" + "-" + strings.ToLower(util.RandomString(5))
	mversion := strings.Split(k8sVersion, ".")

	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range fingerprints {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              name,
		Region:            region,
		Size:              size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-stable",
		},
	}
	tags := []string{"Kubernetes-Cluster", name}

	masterDroplet, _, err := p.createDroplet(client, action, dropletRequest, tags)
	if err != nil {
		return err
	}

	master := strconv.Itoa(masterDroplet.ID)
	m.MasterID = master

	return action.CancellableWaitFor("Kubernetes master launch", 10*time.Minute, 3*time.Second, func() (bool, error) {
		resp, _, serr := client.Droplets.Get(masterDroplet.ID)
		if serr != nil {
			return false, serr
		}

		// Save Master info when ready
		if resp.Status == "active" {
			m.MasterNodes = append(m.MasterNodes, strconv.Itoa(resp.ID))
			m.MasterPrivateIP, _ = resp.PrivateIPv4()
			m.MasterPublicIP, _ = resp.PublicIPv4()
			if serr := p.Core.DB.Save(m); serr != nil {
				return false, serr
			}
		}
		return resp.Status == "active", nil
	})
})
}
