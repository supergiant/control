package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// NOTE despite being named Provider, this is not a full Provider implementation.
// It's really just helper methods shared amongst the providers.

// TODO we should probably make separate interface in core for just this.

type Provider struct {
	Core *core.Core
}

func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	return nil
}

func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	return nil
}

func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
	return nil
}

func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
	return nil
}

func (p *Provider) DeleteNode(m *model.Node, action *core.Action) error {
	return nil
}

func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	service := loadBalancerAsKubernetesService(m)

	if err := p.Core.K8S(m.Kube).CreateResource("api/v1", "Service", m.Namespace, service, service); err != nil {
		return err
	}

	// TODO(stgleb): Context should be inherited from higher level context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	waitDesc := fmt.Sprintf("LoadBalancer %s address", m.Name)
	err := util.WaitFor(ctx, waitDesc, 4*time.Second, func() (bool, error) {
		if getErr := p.Core.K8S(m.Kube).GetResource("api/v1", "Service", m.Namespace, m.Name, service); getErr != nil {
			return false, getErr
		}
		return len(service.Status.LoadBalancer.Ingress) > 0, nil
	})
	if err != nil {
		return err
	}

	return p.Core.DB.Model(m).Update("address", service.Status.LoadBalancer.Ingress[0].Hostname)
}

func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	service := loadBalancerAsKubernetesService(m)
	return p.Core.K8S(m.Kube).UpdateResource("api/v1", "Service", m.Namespace, m.Name, service, service)
}

func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	err := p.Core.K8S(m.Kube).DeleteResource("api/v1", "Service", m.Namespace, m.Name)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	return nil
}

// Private

func loadBalancerAsKubernetesService(m *model.LoadBalancer) *kubernetes.Service {
	var ports []kubernetes.ServicePort
	for lbPort, targetPort := range m.Ports {
		port := kubernetes.ServicePort{
			Port:       lbPort,
			TargetPort: targetPort,
		}
		ports = append(ports, port)
	}

	return &kubernetes.Service{
		Metadata: kubernetes.Metadata{
			Name: m.Name,
		},
		Spec: kubernetes.ServiceSpec{
			Type:     "LoadBalancer",
			Selector: m.Selector,
			Ports:    ports,
		},
	}
}
