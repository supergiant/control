package profile

import (
	"github.com/supergiant/supergiant/pkg/util"
	"sync"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/pkg/errors"
)

const (
	DigitalOceanDefaultProfile = "DigitalOceanDefault"

	OneNodeCluster = "OneNodeK8s"
)

var nodeProfiles *sync.Map
var masterProfiles *sync.Map

//TODO Temporary solution to store profiles. Ideally it should be stored in data storage
func init() {
	nodeProfiles = new(sync.Map)
	nodeProfiles.Store(DigitalOceanDefaultProfile, NodeProfile{
		Name:     DigitalOceanDefaultProfile,
		Provider: util.DigitalOcean,
		OS:       util.UbuntuSlug1604,
		NodeSize: &core.NodeSize{
			Name:     "s-1vcpu-2gb",
			CPUCores: 1,
			RAMGIB:   2,
		},
	})

	masterProfiles = new(sync.Map)
	masterProfiles.Store(OneNodeCluster, ClusterProfile{
		Name:              OneNodeCluster,
		MastersCount:      1,
		NodesCount:        1,
		MasterProfileName: DigitalOceanDefaultProfile,
		NodeProfileName:   DigitalOceanDefaultProfile,
	})
}

//A NodeProfile is the settings of a physical node
type NodeProfile struct {
	Name     string
	OS       string
	Provider util.ProviderName
	Tags     []string
	NodeType string
	NodeSize *core.NodeSize
}

//A cluster profile is the settings of a kubernetes cluster
type ClusterProfile struct {
	Name              string
	MastersCount      int
	NodesCount        int
	MasterProfileName string
	NodeProfileName   string
	K8SVersion        string
}

type NodeProfileService interface {
	GetByName(profileName string) (*NodeProfile, error)
}

type NodeProfiles struct {
}

func (s *NodeProfiles) GetByName(profileName string) (*NodeProfile, error) {
	obj, ok := nodeProfiles.Load(profileName)
	if !ok {
		return nil, errors.Errorf("node profile %s not found", profileName)
	}
	return obj.(*NodeProfile), nil
}

type ClusterProfileService interface {
	GetByName(profileName string) (*ClusterProfile, error)
}

type ClusterProfiles struct {
}

func (s *ClusterProfiles) GetByName(profileName string) (*ClusterProfile, error) {
	obj, ok := nodeProfiles.Load(profileName)
	if !ok {
		return nil, errors.Errorf("cluster profile %s not found", profileName)
	}
	return obj.(*ClusterProfile), nil
}
