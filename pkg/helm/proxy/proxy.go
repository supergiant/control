package proxy

import (
	"fmt"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

const (
	TillerPort = 44134

	TillerConnectionTimeout int64 = 300
)

var (
	tillerPodLabels = labels.Set{"app": "helm", "name": "tiller"}
)

var (
	_ Interface = &Proxy{}
)

// Interface is an interface for tiller proxy.
// It is almost the same as helm interface but without RunReleaseTest function.
type Interface interface {
	ListReleases(opts ...helm.ReleaseListOption) (*rls.ListReleasesResponse, error)
	InstallRelease(chStr, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error)
	InstallReleaseFromChart(chart *chart.Chart, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error)
	DeleteRelease(rlsName string, opts ...helm.DeleteOption) (*rls.UninstallReleaseResponse, error)
	ReleaseStatus(rlsName string, opts ...helm.StatusOption) (*rls.GetReleaseStatusResponse, error)
	UpdateRelease(rlsName, chStr string, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error)
	UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error)
	RollbackRelease(rlsName string, opts ...helm.RollbackOption) (*rls.RollbackReleaseResponse, error)
	ReleaseContent(rlsName string, opts ...helm.ContentOption) (*rls.GetReleaseContentResponse, error)
	ReleaseHistory(rlsName string, opts ...helm.HistoryOption) (*rls.GetHistoryResponse, error)
	GetVersion(opts ...helm.VersionOption) (*rls.GetVersionResponse, error)
	PingTiller() error
}

// Proxy is a wrapper for tiller client for accessing it through kubernetes api.
type Proxy struct {
	kclient         kubernetes.Interface
	restConf        *rest.Config
	tillerNamespace string
}

// New creates a new helm client.
func New(kclient kubernetes.Interface, restConf *rest.Config, tillerNamespace string) (*Proxy, error) {
	return &Proxy{
		kclient:         kclient,
		restConf:        restConf,
		tillerNamespace: tillerNamespace,
	}, nil
}

func (p *Proxy) ListReleases(opts ...helm.ReleaseListOption) (*rls.ListReleasesResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).ListReleases(opts...)
}
func (p *Proxy) InstallRelease(chStr, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).InstallRelease(chStr, namespace, opts...)
}
func (p *Proxy) InstallReleaseFromChart(chart *chart.Chart, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).InstallReleaseFromChart(chart, namespace, opts...)
}
func (p *Proxy) DeleteRelease(rlsName string, opts ...helm.DeleteOption) (*rls.UninstallReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).DeleteRelease(rlsName, opts...)
}
func (p *Proxy) ReleaseStatus(rlsName string, opts ...helm.StatusOption) (*rls.GetReleaseStatusResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).ReleaseStatus(rlsName, opts...)
}
func (p *Proxy) UpdateRelease(rlsName, chStr string, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).UpdateRelease(rlsName, chStr, opts...)
}
func (p *Proxy) UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).UpdateReleaseFromChart(rlsName, chart, opts...)
}
func (p *Proxy) RollbackRelease(rlsName string, opts ...helm.RollbackOption) (*rls.RollbackReleaseResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).RollbackRelease(rlsName, opts...)
}
func (p *Proxy) ReleaseContent(rlsName string, opts ...helm.ContentOption) (*rls.GetReleaseContentResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	return p.helmClient(tun.Local).ReleaseContent(rlsName, opts...)
}
func (p *Proxy) ReleaseHistory(rlsName string, opts ...helm.HistoryOption) (*rls.GetHistoryResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).ReleaseHistory(rlsName, opts...)
}
func (p *Proxy) GetVersion(opts ...helm.VersionOption) (*rls.GetVersionResponse, error) {
	tun, err := p.createTunnel()
	if err != nil {
		return nil, err
	}
	defer tun.close()

	return p.helmClient(tun.Local).GetVersion(opts...)
}

func (p *Proxy) PingTiller() error {
	tun, err := p.createTunnel()
	if err != nil {
		return err
	}
	defer tun.close()

	return p.helmClient(tun.Local).PingTiller()
}

func (p *Proxy) helmClient(port int) helm.Interface {
	return helm.NewClient(helm.Host(fmt.Sprintf("127.0.0.1:%d", port)), helm.ConnectTimeout(TillerConnectionTimeout))
}

// createTunnel creates a tunnel to tiller (like 'kubectl proxy').
func (p *Proxy) createTunnel() (*tunnel, error) {
	tillerPodName, err := getTillerPodName(p.kclient.CoreV1().Pods(p.tillerNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "get tiller pod")
	}

	tun := newTunnel(p.kclient.CoreV1().RESTClient(), p.restConf, p.tillerNamespace, tillerPodName, TillerPort)
	if err := tun.forwardPort(); err != nil {
		return nil, errors.Wrap(err, "setup port forwarding")
	}

	return tun, nil
}

func getTillerPodName(podsClient corev1.PodInterface) (string, error) {
	pod, err := getFirstRunningPod(podsClient, tillerPodLabels.AsSelector())
	if err != nil {
		return "", err
	}

	return pod.ObjectMeta.GetName(), nil
}

func getFirstRunningPod(podsClient corev1.PodInterface, selector labels.Selector) (*apiv1.Pod, error) {
	pods, err := podsClient.List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, errors.New("could not find tiller")
	}

	for _, p := range pods.Items {
		if isPodReady(&p) {
			return &p, nil
		}
	}

	return nil, errors.New("could not find a ready tiller pod")
}
