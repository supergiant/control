package profile

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/provider"
)

const (
	//DigitalOceanDefaultProfile is default node profile for digital ocean
	DigitalOceanDefaultProfile = "DigitalOceanDefault"
	//OneNodeCluster is default kubernetes one master one node kube profile
	OneNodeCluster = "OneNodeK8s"
)

var nodeProfiles *sync.Map
var masterProfiles *sync.Map

//TODO Temporary solution to store profiles. Ideally it should be stored in data storage
func init() {
	nodeProfiles = new(sync.Map)
	nodeProfiles.Store(DigitalOceanDefaultProfile, NodeProfile{
		Name:     DigitalOceanDefaultProfile,
		Provider: provider.DigitalOcean,
		OS:       util.UbuntuSlug1604,
		NodeSize: &core.NodeSize{
			Name:     "s-1vcpu-2gb",
			CPUCores: 1,
			RAMGIB:   2,
		},
	})

	masterProfiles = new(sync.Map)
	masterProfiles.Store(OneNodeCluster, KubeProfile{
		Name:              OneNodeCluster,
		MastersCount:      1,
		NodesCount:        1,
		MasterProfileName: DigitalOceanDefaultProfile,
		NodeProfileName:   DigitalOceanDefaultProfile,
	})
}

//NodeProfile is the settings of a physical node
type NodeProfile struct {
	Name     string
	OS       string
	Provider provider.Name
	Tags     []string
	NodeType string
	NodeSize *core.NodeSize
}

//KubeProfile is the settings for a k8s cluster
type KubeProfile struct {
	Name              string
	MastersCount      int
	NodesCount        int
	MasterProfileName string
	NodeProfileName   string
	K8SVersion        string
}

//NodeProfileService provides operations over node profiles
type NodeProfileService interface {
	GetByName(profileName string) (*NodeProfile, error)
}

//NodeProfiles is an implementation of NodeProfileService
type NodeProfiles struct {
}

//GetByName retrieves a node profile from the storage by it's unique name
func (s *NodeProfiles) GetByName(profileName string) (*NodeProfile, error) {
	obj, ok := nodeProfiles.Load(profileName)
	if !ok {
		return nil, errors.Errorf("node profile %s not found", profileName)
	}
	return obj.(*NodeProfile), nil
}

//KubeProfileService provides operations over cluster profiles
type KubeProfileService interface {
	GetByName(profileName string) (*KubeProfile, error)
}

//KubeProfiles is an implementation of KubeProfileService
type KubeProfiles struct {
}

//GetByName retrieves a kube profile from the storage by it's unique name
func (s *KubeProfiles) GetByName(profileName string) (*KubeProfile, error) {
	obj, ok := nodeProfiles.Load(profileName)
	if !ok {
		return nil, errors.Errorf("cluster profile %s not found", profileName)
	}
	return obj.(*KubeProfile), nil
}
