package fake_digitalocean_provider

import "github.com/digitalocean/godo"

type Storage struct {
	ListVolumesFn  func(*godo.ListOptions) ([]godo.Volume, *godo.Response, error)
	GetVolumeFn    func(string) (*godo.Volume, *godo.Response, error)
	CreateVolumeFn func(*godo.VolumeCreateRequest) (*godo.Volume, *godo.Response, error)
	DeleteVolumeFn func(string) (*godo.Response, error)
}

func (f *Storage) ListVolumes(l *godo.ListOptions) ([]godo.Volume, *godo.Response, error) {
	if f.ListVolumesFn == nil {
		return nil, nil, nil
	}
	return f.ListVolumesFn(l)
}

func (f *Storage) GetVolume(s string) (*godo.Volume, *godo.Response, error) {
	if f.GetVolumeFn == nil {
		return nil, nil, nil
	}
	return f.GetVolumeFn(s)
}

func (f *Storage) CreateVolume(r *godo.VolumeCreateRequest) (*godo.Volume, *godo.Response, error) {
	if f.CreateVolumeFn == nil {
		return nil, nil, nil
	}
	return f.CreateVolumeFn(r)
}

func (f *Storage) DeleteVolume(s string) (*godo.Response, error) {
	if f.DeleteVolumeFn == nil {
		return nil, nil
	}
	return f.DeleteVolumeFn(s)
}
