package kube

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/storage"
	"bytes"
	"github.com/sirupsen/logrus"
)

const prefix = "/kube/"

type Service struct {
	provider   core.Provider
	repository storage.Interface
}

func (s *Service) CreateKube(ctx context.Context, kube *Kube) (*Kube, error) {
	logrus.Debug("kube.Service.CreateKube start")
	//TODO VALIDATION
	raw, err := json.Marshal(kube)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = s.repository.Put(ctx, prefix, kube.KubeName, raw)
	if err != nil {
		return nil, err
	}

	// TODO Add provisioning/provider
	// TODO Add async task handling
	//s.provider.CreateKube()
	logrus.Debug("kube.Service.CreateKube end")
	return nil, nil
}

func (s *Service) ListAll(ctx context.Context) ([]Kube, error) {
	logrus.Debug("kube.Service.ListAll start")

	rawKubes, err := s.repository.GetAll(ctx, prefix)
	if err != nil {
		return nil, err
	}
	kubes := make([]Kube, 0)
	for _, v := range rawKubes {
		k := &Kube{}
		err = json.NewDecoder(bytes.NewReader(v)).Decode(k)
		if err != nil {
			logrus.Warning("corrupted data: can't read kube struct from storage!")
			continue
		}
		kubes = append(kubes, *k)
	}

	logrus.Debug("kube.Service.ListAll end")
	return kubes, nil
}

func (s *Service) Get(ctx context.Context, kubeName string) (*Kube, error) {
	logrus.Debug("kube.Service.Get start")
	rawJSON, err := s.repository.Get(ctx, prefix, kubeName)
	if err != nil {
		return nil, err
	}
	if rawJSON == nil {
		return nil, nil
	}
	k := &Kube{}
	err = json.NewDecoder(bytes.NewReader(rawJSON)).Decode(k)
	if err != nil {
		return nil, err
	}
	logrus.Debug("kube.Service.Get end")
	return k, nil
}

func (s *Service) Update(ctx context.Context, kube *Kube) (*Kube, error) {
	logrus.Debug("kube.Service.Update start")
	panic("Implement me")
	logrus.Debug("kube.Service.Update end")
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, kubeName string) error {
	logrus.Debug("kube.Service.Delete start")
	//TODO Delete load balancers, nodes, cluster
	panic("Implement me")
	logrus.Debug("kube.Service.Delete end")
	return nil
}
