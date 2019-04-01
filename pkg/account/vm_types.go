package account

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
)

//go:generate go run awsvms/main.go

type MachineTypes struct {
	rmu     sync.RWMutex
	regions []string

	vmu       sync.RWMutex
	regionVMs map[string][]string

	tmu     sync.RWMutex
	vmTypes map[string]VMType
}

type VMType struct {
	Name      string
	VCPU      string
	MemoryGiB string
	GPU       string
}

// Regions return a list of supported regions.
func (s MachineTypes) Regions() []string {
	s.rmu.RLock()
	defer s.rmu.RUnlock()

	return s.regions
}

// Region return a list of supported vm types in the region. Used by RegionFinders.
func (s MachineTypes) RegionTypes(region string) ([]string, error) {
	s.vmu.RLock()
	defer s.vmu.RUnlock()

	// TODO: check reginVMs is nil?
	if _, ok := s.regionVMs[region]; !ok {
		return nil, errors.Wrapf(sgerrors.ErrRawError, "unknown region: %s", region)
	}
	return s.regionVMs[region], nil
}

// Sizes provides virtual machine parameters (cpu/ram/gpu). Used by RegionFinders.
func (s MachineTypes) Sizes() map[string]VMType {
	s.tmu.RLock()
	defer s.tmu.RUnlock()

	return s.vmTypes
}
