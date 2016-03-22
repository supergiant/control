package core

import (
	"guber"
	"strconv"
	"strings"
)

// NOTE the word Blueprint is used for Volumes and Containers, since they are
// both "definitions" that create "instances" of the real thing

type Blueprint struct {
	Volumes                []*VolumeBlueprint    `json:"volumes"`
	Containers             []*ContainerBlueprint `json:"containers"`
	TerminationGracePeriod int                   `json:"termination_grace_period"`
}

// Volume
//==============================================================================
type VolumeBlueprint struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

// Container
//==============================================================================
type ContainerBlueprint struct {
	Image  string              `json:"image"`
	Ports  []*Port             `json:"ports"`
	Env    []*EnvVar           `json:"env"`
	CPU    *ResourceAllocation `json:"cpu"`
	RAM    *ResourceAllocation `json:"ram"`
	Mounts []*Mount            `json:"mounts"`
}

type ResourceAllocation struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func (m *ContainerBlueprint) kubeVolumeMounts() (volMounts []*guber.VolumeMount) {
	for _, mount := range m.Mounts {
		volMounts = append(volMounts, mount.asKubeVolumeMount())
	}
	return volMounts
}

func (m *ContainerBlueprint) kubeContainerPorts() (cPorts []*guber.ContainerPort) {
	for _, port := range m.Ports {
		cPorts = append(cPorts, port.asKubeContainerPort())
	}
	return cPorts
}

func (m *ContainerBlueprint) interpolatedEnvVars(instance *Instance) (envVars []*guber.EnvVar) {
	for _, envVar := range m.Env {
		envVars = append(envVars, envVar.asKubeEnvVar(instance))
	}
	return envVars
}

func (m *ContainerBlueprint) ImageRepoName() string {
	return strings.Split(m.Image, "/")[0]
}

func (m *ContainerBlueprint) AsKubeContainer(instance *Instance) *guber.Container { // NOTE how instance must be passed here
	return &guber.Container{
		Name:  strconv.Itoa(instance.ID), // TODO this is not good naming
		Image: m.Image,
		Env:   m.interpolatedEnvVars(instance),
		Resources: &guber.Resources{
			Requests: &guber.ResourceValues{
				Memory: BytesFromMiB(m.RAM.Min).ToKubeMebibytes(),
				CPU:    (&CoresValue{m.CPU.Min}).ToKubeMillicores(),
			},
			Limits: &guber.ResourceValues{
				Memory: BytesFromMiB(m.RAM.Max).ToKubeMebibytes(),
				CPU:    (&CoresValue{m.CPU.Max}).ToKubeMillicores(),
			},
		},
		VolumeMounts: m.kubeVolumeMounts(),
		Ports:        m.kubeContainerPorts(),
	}
}

// EnvVar
//==============================================================================
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"` // this may be templated, "something_{{ instance_id }}"
}

func (m *EnvVar) interpolatedValue(instance *Instance) string {
	r := strings.NewReplacer("{{ instance_id }}", strconv.Itoa(instance.ID),
		"{{ other_stuff }}", "TODO") // TODO
	return r.Replace(m.Value)
}

func (m *EnvVar) asKubeEnvVar(instance *Instance) *guber.EnvVar {
	return &guber.EnvVar{
		Name:  m.Name,
		Value: m.interpolatedValue(instance),
	}
}

// Mount
//==============================================================================
type Mount struct {
	Volume string `json:"volume"`
	Path   string `json:"path"`
}

func (m *Mount) asKubeVolumeMount() *guber.VolumeMount {
	return &guber.VolumeMount{
		Name:      m.Volume,
		MountPath: m.Path,
	}
}

// Port
//==============================================================================
type Port struct {
	Protocol string `json:"protocol"`
	Number   int    `json:"number"`
	Public   bool   `json:"public"`
}

// TODO not sure if ports should be named this way
func (m *Port) name() string {
	return strconv.Itoa(m.Number)
}

func (m *Port) asKubeContainerPort() *guber.ContainerPort {
	return &guber.ContainerPort{
		ContainerPort: m.Number,
	}
}

func (m *Port) AsKubeServicePort() *guber.ServicePort {
	return &guber.ServicePort{
		Name:     m.name(),
		Port:     m.Number,
		Protocol: "TCP", // this is default; only other option is UDP
	}
}
