package digitalocean

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/template"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateKube creates a new DO kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {

	if m.SSHPubKey == "" {
		m.SSHPubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== schacon@mylaptop.local"
	}

	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	client := p.Client(m)

	// Default master count to 1
	if m.KubeMasterCount == 0 {
		m.KubeMasterCount = 1
	}

	m.CustomFiles = fmt.Sprintf(`
  - path: "/root/.config/doctl/config.yaml"
    permissions: "0600"
    owner: "root"
    content: |
      access-token: %s
      output: text
  - path: "/etc/kubernetes/volumeplugins/supergiant.io~digitalocean/digitalocean"
    permissions: "0755"
    owner: "root"
    content: |
      #!/bin/bash
      # Required Flex Volume Options.
      #{
      #  "volumeID": "bar",
      #  "name": "foo"
      #}
      ## CoreOS Kube-Wrapper setup
      mount -o remount,ro /sys/fs/selinux >/dev/null 2>&1
      export PATH=$PATH:/usr/local/sbin/
      apt-get update >/dev/null 2>&1
      apt-get install -y curl wget jq >/dev/null 2>&1
      mkdir -p /opt/bin/ >/dev/null 2>&1
      wget https://github.com/digitalocean/doctl/releases/download/v1.4.0/doctl-1.4.0-linux-amd64.tar.gz >/dev/null 2>&1
      tar xf doctl-1.4.0-linux-amd64.tar.gz >/dev/null 2>&1
      mv ./doctl /opt/bin/ >/dev/null 2>&1
      mount -o remount,rw /sys/fs/selinux >/dev/null 2>&1
      export PATH=$PATH:.
      # Who am i?
      # Where am i?
      PUBLICIP=$(wget http://ipinfo.io/ip -qO -)
      REGION=$(/opt/bin/doctl compute droplet list --config /root/.config/doctl/config.yaml| grep ${PUBLICIP} | awk '{print $7}')
      DROPLET_ID=$(/opt/bin/doctl compute droplet list --config /root/.config/doctl/config.yaml| grep ${PUBLICIP} | awk '{print $1}')

      usage() {
      	err "Invalid usage. Usage: "
      	err "\t$0 init"
      	err "\t$0 attach <json params>"
      	err "\t$0 detach <mount device>"
      	err "\t$0 mount <mount dir> <mount device> <json params>"
      	err "\t$0 unmount <mount dir>"
      	exit 1
      }

      err() {
      	echo -ne $* 1>&2
      }

      log() {
      	echo -ne $* >&1
      }

      ismounted() {
      	MOUNT=$(findmnt -n ${MNTPATH} 2>/dev/null | cut -d' ' -f1)
      	if [ "${MOUNT}" == "${MNTPATH}" ]; then
      		echo "1"
      	else
      		echo "0"
      	fi
      }

      attach() {
      	VOLUMEID=$(echo $1 | jq -r '.volumeID')
      	VOLUMENAME=$(echo $1 | jq -r '.name')
        /opt/bin/doctl compute volume-action attach $VOLUMEID $DROPLET_ID --config /root/.config/doctl/config.yaml >/dev/null 2>&1
        # Find the new volume.
      	DEVNAME="/dev/disk/by-id/scsi-0DO_Volume_${VOLUMENAME}"
      	# Wait for attach
      	NEXT_WAIT_TIME=1
        until ls -l $DEVNAME >/dev/null 2>&1 || [ $NEXT_WAIT_TIME -eq 4 ]; do
         sleep $(( NEXT_WAIT_TIME++ ))
        done
      	#Record the actual device name.
      	DVSHRTNAME=$(ls -l /dev/disk/by-id | grep ${VOLUMENAME} | awk '{print $11}' | sed 's/\.\.\///g' | sed '/^\s*$/d')
      	DMDEV="/dev/${DVSHRTNAME}"
      	# Error check.
      	if [ ! -b "${DMDEV}" ]; then
      		err "{\"status\": \"Failure\", \"message\": \"Volume ${DMDEV} does not exist\"}"
      		exit 1
      	fi
      	log "{\"status\": \"Success\", \"device\":\"${DMDEV}\"}"
      	exit 0
      }

      detach() {
      	# This is nasty, I would prefer to use doctl for detach as well... but it appears that it is bugged.
      	# I will update this when a new version of doctl releases. For now raw api.
      	TOKEN=$(cat ~/.config/doctl/config.yaml | grep access-token | awk '{print $2}')
      	SRTDEVNAME=$(echo $1 | sed 's/\/dev\///')
      	VOLNAME=$(ls -l /dev/disk/by-id | grep ${SRTDEVNAME} | awk '{print $9}' | sed 's/scsi-0DO_Volume_//')
      	curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer ${TOKEN}" -d "{\"type\": \"detach\", \"droplet_id\": \"${DROPLET_ID}\", \"volume_name\": \"${VOLNAME}\", \"region\": \"nyc1\"}" "https://api.digitalocean.com/v2/volumes/actions" >/dev/null 2>&1
      	if [ -b "$1" ]; then
      		log "{\"status\": \"Success\"}"
      		exit 0
      	fi
      	exit 1
      }

      domount() {
      	MNTPATH=$1
      	DMDEV=$2
      	FSTYPE=$(echo $3|jq -r '.["kubernetes.io/fsType"]')
      	if [ ! -b "${DMDEV}" ]; then
      		err "{\"status\": \"Failure\", \"message\": \"${DMDEV} does not exist\"}"
      		exit 1
      	fi
      	if [ $(ismounted) -eq 1 ] ; then
      		log "{\"status\": \"Success\"}"
      		exit 0
      	fi
      	VOLFSTYPE=$(blkid -o udev ${DMDEV} 2>/dev/null|grep "ID_FS_TYPE"|cut -d"=" -f2)
      	if [ "${VOLFSTYPE}" == "" ]; then
      		mkfs -t ${FSTYPE} ${DMDEV} >/dev/null 2>&1
      		if [ $? -ne 0 ]; then
      			err "{ \"status\": \"Failure\", \"message\": \"Failed to create fs ${FSTYPE} on device ${DMDEV}\"}"
      			exit 1
      		fi
      	fi
      	mkdir -p ${MNTPATH} &> /dev/null
      	mount ${DMDEV} ${MNTPATH} &> /dev/null
      	if [ $? -ne 0 ]; then
      		err "{ \"status\": \"Failure\", \"message\": \"Failed to mount device ${DMDEV} at ${MNTPATH}\"}"
      		exit 1
      	fi
      	log "{\"status\": \"Success\"}"
      	exit 0
      }

      unmount() {
      	MNTPATH=$1
      	if [ $(ismounted) -eq 0 ] ; then
      		log "{\"status\": \"Success\"}"
      		exit 0
      	fi
      	umount ${MNTPATH} &> /dev/null
      	if [ $? -ne 0 ]; then
      		err "{ \"status\": \"Failed\", \"message\": \"Failed to unmount volume at ${MNTPATH}\"}"
      		exit 1
      	fi
      	rmdir ${MNTPATH} &> /dev/null
      	log "{\"status\": \"Success\"}"
      	exit 0
      }
      op=$1
      if [ "$op" = "init" ]; then
        log "{\"status\": \"Success\"}"
        exit 0
      fi
      if [ $# -lt 2 ]; then
      	usage
      fi
      shift
      case "$op" in
      	attach)
      		attach $*
      		;;
      	detach)
      		detach $*
      		;;
      	mount)
      		domount $*
      		;;
      	unmount)
      		unmount $*
      		;;
      	*)
      		usage
      esac

      exit 1`, m.CloudAccount.Credentials["token"])

	// provision an etcd token
	url, err := etcdToken(strconv.Itoa(m.KubeMasterCount))
	if err != nil {
		return err
	}

	err = p.Core.DB.Save(m)
	if err != nil {
		return err
	}
	// save the token
	m.ETCDDiscoveryURL = url

	procedure.AddStep("creating global tags for Kube", func() error {
		// These are created once, and then attached by name to created resource
		globalTags := []string{
			"Kubernetes-Cluster",
			m.Name,
			m.Name + "-master",
			m.Name + "-minion",
		}
		for _, tag := range globalTags {
			createInput := &godo.TagCreateRequest{
				Name: tag,
			}
			if _, _, err := client.Tags.Create(createInput); err != nil {
				// TODO
				p.Core.Log.Warnf("Failed to create Digital Ocean tag '%s': %s", tag, err)
			}
		}
		return nil
	})

	for i := 1; i <= m.KubeMasterCount; i++ {
		// Create master(s)
		count := strconv.Itoa(i)

		procedure.AddStep("Creating Kubernetes Master Node "+count+"...", func() error {

			// Master name
			name := m.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))

			m.MasterName = name

			// Build template
			mversion := strings.Split(m.KubernetesVersion, ".")
			masterFileName := fmt.Sprintf("config/providers/common/%s.%s/master.yaml)", mversion[0], mversion[1])
			masterTemplate := template.Templates[masterFileName]

			var masterUserdata bytes.Buffer
			if err = masterTemplate.Execute(&masterUserdata, m); err != nil {
				return err
			}

			var fingers []godo.DropletCreateSSHKey
			for _, ssh := range m.DigitalOceanConfig.SSHKeyFingerprint {
				fingers = append(fingers, godo.DropletCreateSSHKey{
					Fingerprint: ssh,
				})
			}

			dropletRequest := &godo.DropletCreateRequest{
				Name:              name,
				Region:            m.DigitalOceanConfig.Region,
				Size:              m.MasterNodeSize,
				PrivateNetworking: true,
				UserData:          string(masterUserdata.Bytes()),
				SSHKeys:           fingers,
				Image: godo.DropletCreateImage{
					Slug: "coreos-stable",
				},
			}
			tags := []string{"Kubernetes-Cluster", m.Name, name}

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

	procedure.AddStep("building Kubernetes minion", func() error {
		// Load Nodes to see if we've already created a minion
		// TODO -- I think we can get rid of a lot of this do-unless behavior if we
		// modify Procedure to save progess on Action (which is easy to implement).
		if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		}

		node := &model.Node{
			KubeName: m.Name,
			Kube:     m,
			Size:     m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})

	// TODO repeated in provider_aws.go
	procedure.AddStep("waiting for Kubernetes", func() error {
		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, 3*time.Second, func() (bool, error) {
			k8s := p.Core.K8S(m)
			k8sNodes, err := k8s.ListNodes("")
			if err != nil {
				return false, nil
			}
			return len(k8sNodes) > 0, nil
		})
	})

	return procedure.Run()
}
