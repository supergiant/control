package fake_digitalocean_provider

import "github.com/digitalocean/godo"

type Droplets struct {
	ListFn           func(*godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	ListByTagFn      func(string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	GetFn            func(int) (*godo.Droplet, *godo.Response, error)
	CreateFn         func(*godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
	CreateMultipleFn func(*godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error)
	DeleteFn         func(int) (*godo.Response, error)
	DeleteByTagFn    func(string) (*godo.Response, error)
	KernelsFn        func(int, *godo.ListOptions) ([]godo.Kernel, *godo.Response, error)
	SnapshotsFn      func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	BackupsFn        func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	ActionsFn        func(int, *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	NeighborsFn      func(int) ([]godo.Droplet, *godo.Response, error)
}

func (f *Droplets) List(l *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	if f.ListFn == nil {
		return nil, nil, nil
	}
	return f.ListFn(l)
}

func (f *Droplets) ListByTag(s string, l *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	if f.ListByTagFn == nil {
		return nil, nil, nil
	}
	return f.ListByTagFn(s, l)
}

func (f *Droplets) Get(id int) (*godo.Droplet, *godo.Response, error) {
	if f.GetFn == nil {
		return nil, nil, nil
	}
	return f.GetFn(id)
}

func (f *Droplets) Create(r *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
	if f.CreateFn == nil {
		return nil, nil, nil
	}
	return f.CreateFn(r)
}

func (f *Droplets) CreateMultiple(r *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error) {
	if f.CreateMultipleFn == nil {
		return nil, nil, nil
	}
	return f.CreateMultipleFn(r)
}

func (f *Droplets) Delete(id int) (*godo.Response, error) {
	if f.DeleteFn == nil {
		return nil, nil
	}
	return f.DeleteFn(id)
}

func (f *Droplets) DeleteByTag(s string) (*godo.Response, error) {
	if f.DeleteByTagFn == nil {
		return nil, nil
	}
	return f.DeleteByTagFn(s)
}

func (f *Droplets) Kernels(id int, l *godo.ListOptions) ([]godo.Kernel, *godo.Response, error) {
	if f.KernelsFn == nil {
		return nil, nil, nil
	}
	return f.KernelsFn(id, l)
}

func (f *Droplets) Snapshots(id int, l *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	if f.SnapshotsFn == nil {
		return nil, nil, nil
	}
	return f.SnapshotsFn(id, l)
}

func (f *Droplets) Backups(id int, l *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	if f.BackupsFn == nil {
		return nil, nil, nil
	}
	return f.BackupsFn(id, l)
}

func (f *Droplets) Actions(id int, l *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	if f.ActionsFn == nil {
		return nil, nil, nil
	}
	return f.ActionsFn(id, l)
}

func (f *Droplets) Neighbors(id int) ([]godo.Droplet, *godo.Response, error) {
	if f.NeighborsFn == nil {
		return nil, nil, nil
	}
	return f.NeighborsFn(id)
}
