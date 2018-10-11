package kube

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/technosophos/moniker"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/timeconv"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/sghelm"
	"github.com/supergiant/supergiant/pkg/sghelm/proxy"
	"github.com/supergiant/supergiant/pkg/storage"
)

const (
	DefaultStoragePrefix = "/supergiant/kube/"

	releaseInstallTimeout = 300
)

// Interface represents an interface for a kube service.
type Interface interface {
	Create(ctx context.Context, k *model.Kube) error
	Get(ctx context.Context, name string) (*model.Kube, error)
	ListAll(ctx context.Context) ([]model.Kube, error)
	Delete(ctx context.Context, name string) error
	ListKubeResources(ctx context.Context, kname string) ([]byte, error)
	GetKubeResources(ctx context.Context, kname, resource, ns, name string) ([]byte, error)
	GetCerts(ctx context.Context, kname, cname string) (*Bundle, error)
}

// Service manages kubernetes clusters.
type Service struct {
	discoveryClientFn func(k *model.Kube) (*discovery.DiscoveryClient, error)
	clientForGroupFn  func(k *model.Kube, gv schema.GroupVersion) (rest.Interface, error)

	prefix  string
	storage storage.Interface

	helmRepos sghelm.Servicer
}

// NewService constructs a Service.
func NewService(prefix string, s storage.Interface, helmSvc sghelm.Servicer) *Service {
	return &Service{
		clientForGroupFn:  restClientForGroupVersion,
		discoveryClientFn: discoveryClient,
		helmRepos:         helmSvc,
		prefix:            prefix,
		storage:           s,
	}
}

// Create and stores a kube in the provided storage.
func (s *Service) Create(ctx context.Context, k *model.Kube) error {
	raw, err := json.Marshal(k)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	err = s.storage.Put(ctx, s.prefix, k.Name, raw)
	if err != nil {
		return errors.Wrap(err, "storage: put")
	}

	return nil
}

// Get returns a kube with a specified name.
func (s *Service) Get(ctx context.Context, name string) (*model.Kube, error) {
	raw, err := s.storage.Get(ctx, s.prefix, name)
	if err != nil {
		return nil, errors.Wrap(err, "storage: get")
	}
	if raw == nil {
		return nil, sgerrors.ErrNotFound
	}

	k := &model.Kube{}
	if err = json.Unmarshal(raw, k); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return k, nil
}

// ListAll returns all kubes.
func (s *Service) ListAll(ctx context.Context) ([]model.Kube, error) {
	rawKubes, err := s.storage.GetAll(ctx, s.prefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage: getAll")
	}

	kubes := make([]model.Kube, len(rawKubes))
	for i, v := range rawKubes {
		k := model.Kube{}
		if err = json.Unmarshal(v, &k); err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}
		kubes[i] = k
	}

	return kubes, nil
}

// Delete deletes a kube with a specified name.
func (s *Service) Delete(ctx context.Context, name string) error {
	return s.storage.Delete(ctx, s.prefix, name)
}

// ListKubeResources returns raw representation of the supported kubernetes resources.
func (s *Service) ListKubeResources(ctx context.Context, kname string) ([]byte, error) {
	kube, err := s.Get(ctx, kname)
	if err != nil {
		return nil, errors.Wrap(err, "storage: get")
	}

	resourcesInfo, err := s.resourcesGroupInfo(kube)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(resourcesInfo)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	return raw, nil
}

// GetKubeResources returns raw representation of the kubernetes resources.
func (s *Service) GetKubeResources(ctx context.Context, kname, resource, ns, name string) ([]byte, error) {
	kube, err := s.Get(ctx, kname)
	if err != nil {
		return nil, errors.Wrap(err, "storage: get")
	}

	resourcesInfo, err := s.resourcesGroupInfo(kube)
	if err != nil {
		return nil, err
	}

	gv, ok := resourcesInfo[resource]
	if !ok {
		return nil, sgerrors.ErrNotFound
	}

	client, err := s.clientForGroupFn(kube, gv)
	if err != nil {
		return nil, errors.Wrap(err, "get kube client")
	}

	raw, err := client.Get().Resource(resource).Namespace(ns).Name(name).DoRaw()
	if err != nil {
		return nil, errors.Wrap(err, "get resources")
	}

	return raw, nil
}

// GetCerts returns a keys bundle for provided component name.
func (s *Service) GetCerts(ctx context.Context, kname, cname string) (*Bundle, error) {
	kube, err := s.Get(ctx, kname)
	if err != nil {
		return nil, err
	}

	r, err := ssh.NewRunner(ssh.Config{
		User: kube.SshUser,
		Key:  kube.SshPublicKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup runner")
	}

	certs, err := NewCerts(DefaultCertsPath, r)
	if err != nil {
		return nil, errors.Wrap(err, "setup certs getter")
	}

	b, err := certs.BundleFor(ctx, cname)
	if err != nil {
		return nil, errors.Wrap(err, "get keys bundle")
	}

	return b, nil
}

func (s *Service) InstallChart(ctx context.Context, kname string, rls *ReleaseInput) (*model.ReleaseInfo, error) {
	kube, err := s.Get(ctx, kname)
	if err != nil {
		return nil, err
	}

	chrt, err := s.helmRepos.GetChart(ctx, rls.RepoName, rls.ChartName, rls.ChartVersion)
	if err != nil {
		return nil, err
	}

	restConf, err := NewConfigFor(kube)
	if err != nil {
		return nil, err
	}
	coreV1Client, err := corev1.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}

	kprx, err := proxy.New(coreV1Client, restConf, "")
	if err != nil {
		return nil, err
	}

	rr, err := kprx.InstallReleaseFromChart(
		chrt,
		rls.Namespace,
		helm.ReleaseName(ensureReleaseName(rls.Name)),
		helm.ValueOverrides(rls.Values),
		helm.InstallWait(false),
		helm.InstallTimeout(releaseInstallTimeout),
	)

	return toReleaseInfo(rr.Release), err
}

func (s *Service) resourcesGroupInfo(kube *model.Kube) (map[string]schema.GroupVersion, error) {
	client, err := s.discoveryClientFn(kube)
	if err != nil {
		return nil, errors.Wrap(err, "get discovery client")
	}

	apiResourceLists, err := client.ServerResources()
	if err != nil {
		return nil, errors.Wrap(err, "get resources")
	}

	resourcesGroupInfo := map[string]schema.GroupVersion{}
	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			if _, ok := resourcesGroupInfo[apiResource.Kind]; !ok {
				resourcesGroupInfo[apiResource.Name] = gv
			}
		}
	}

	return resourcesGroupInfo, nil
}

func ensureReleaseName(name string) string {
	if strings.TrimSpace(name) == "" {
		return moniker.New().NameSep("-")
	}
	return name
}

func toReleaseInfo(rls *release.Release) *model.ReleaseInfo {
	if rls == nil {
		return nil
	}
	return &model.ReleaseInfo{
		Name:         rls.GetName(),
		Namespace:    rls.GetNamespace(),
		Version:      rls.GetVersion(),
		CreatedAt:    timeconv.String(rls.GetInfo().GetFirstDeployed()),
		LastDeployed: timeconv.String(rls.GetInfo().GetLastDeployed()),
		Chart:        rls.GetChart().Metadata.Name,
		ChartVersion: rls.GetChart().Metadata.Version,
		Status:       rls.GetInfo().Status.String(),
	}
}
