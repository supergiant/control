package kube

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

const prefix = "/supergiant/kube/"

// Service manages kubernetes clusters.
type Service struct {
	discoveryClientFn func(k *Kube) (*discovery.DiscoveryClient, error)
	clientForGroupFn  func(k *Kube, gv schema.GroupVersion) (rest.Interface, error)

	storage storage.Interface
}

// NewService constructs a Service.
func NewService(s storage.Interface) *Service {
	return &Service{
		clientForGroupFn:  restClientForGroupVersion,
		discoveryClientFn: discoveryClient,
		storage:           s,
	}
}

// Create stores a kube in the provided storage.
func (s *Service) CreateKube(ctx context.Context, k *Kube) (*Kube, error) {
	k.ID = getHash(k.APIHost + k.APIPort + k.Auth.Username)
	if k.ID == strings.TrimSpace(k.ID) {
		return nil, ErrInvalidID
	}

	raw, err := json.Marshal(k)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	err = s.storage.Put(ctx, prefix, k.ID, raw)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	return k, nil
}

// Get returns a kube with a specified id.
func (s *Service) Get(ctx context.Context, kubeID string) (*Kube, error) {
	raw, err := s.storage.Get(ctx, prefix, kubeID)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}
	if raw == nil {
		return nil, nil
	}

	k := &Kube{}
	err = json.NewDecoder(bytes.NewReader(raw)).Decode(k)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return k, nil
}

// Get returns all kubes.
func (s *Service) ListAll(ctx context.Context) ([]Kube, error) {
	rawKubes, err := s.storage.GetAll(ctx, prefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	kubes := make([]Kube, len(rawKubes))
	for i, v := range rawKubes {
		k := &Kube{}
		err = json.NewDecoder(bytes.NewReader(v)).Decode(k)
		if err != nil {
			logrus.Warning("corrupted data: can't read kube struct from storage!")
			continue
		}
		kubes[i] = *k
	}

	return kubes, nil
}

// Get deletes a kube with a specified id.
func (s *Service) Delete(ctx context.Context, kubeID string) error {
	return s.storage.Delete(ctx, prefix, kubeID)
}

// ListKubeResources returns raw representation of the supported kubernetes resources.
func (s *Service) ListKubeResources(ctx context.Context, kubeID string) ([]byte, error) {
	kube, err := s.Get(ctx, kubeID)
	if err != nil {
		return nil, errors.Wrap(err, "get kube")
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
func (s *Service) GetKubeResources(ctx context.Context, kubeID, resource, ns, name string) ([]byte, error) {
	kube, err := s.Get(ctx, kubeID)
	if err != nil {
		return nil, errors.Wrap(err, "get kube")
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

func getHash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
