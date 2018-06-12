package kube

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

const DefaultStoragePrefix = "/kube/"

// Service manages kubernetes clusters.
type Service struct {
	discoveryClientFn func(k *Kube) (*discovery.DiscoveryClient, error)
	clientForGroupFn  func(k *Kube, gv schema.GroupVersion) (rest.Interface, error)

	prefix  string
	storage storage.Interface
}

// NewService constructs a Service.
func NewService(prefix string, s storage.Interface) *Service {
	return &Service{
		clientForGroupFn:  restClientForGroupVersion,
		discoveryClientFn: discoveryClient,
		prefix:            prefix,
		storage:           s,
	}
}

// Create stores a kube in the provided storage.
func (s *Service) Create(ctx context.Context, k *Kube) error {
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
func (s *Service) Get(ctx context.Context, name string) (*Kube, error) {
	raw, err := s.storage.Get(ctx, s.prefix, name)
	if err != nil {
		return nil, errors.Wrap(err, "storage: get")
	}
	if raw == nil {
		return nil, sgerrors.ErrNotFound
	}

	k := &Kube{}
	if err = json.Unmarshal(raw, k); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return k, nil
}

// ListAll returns all kubes.
func (s *Service) ListAll(ctx context.Context) ([]Kube, error) {
	rawKubes, err := s.storage.GetAll(ctx, s.prefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage: getAll")
	}

	kubes := make([]Kube, len(rawKubes))
	for i, v := range rawKubes {
		k := Kube{}
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
func (s *Service) ListKubeResources(ctx context.Context, name string) ([]byte, error) {
	kube, err := s.Get(ctx, name)
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
func (s *Service) GetKubeResources(ctx context.Context, kubeName, resource, ns, name string) ([]byte, error) {
	kube, err := s.Get(ctx, kubeName)
	if err != nil {
		return nil, errors.Wrap(err, "storage: get")
	}

	resourcesInfo, err := s.resourcesGroupInfo(kube)
	if err != nil {
		return nil, err
	}

	gv, ok := resourcesInfo[resource]
	if !ok {
		return nil, ErrNotFound
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

func (s *Service) resourcesGroupInfo(kube *Kube) (map[string]schema.GroupVersion, error) {
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
