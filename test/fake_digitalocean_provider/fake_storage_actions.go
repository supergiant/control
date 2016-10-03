package fake_digitalocean_provider

import "github.com/digitalocean/godo"

type StorageActions struct {
	AttachFn func(volumeID string, dropletID int) (*godo.Action, *godo.Response, error)
	DetachFn func(volumeID string) (*godo.Action, *godo.Response, error)
	GetFn    func(volumeID string, actionID int) (*godo.Action, *godo.Response, error)
	ListFn   func(volumeID string, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	ResizeFn func(volumeID string, sizeGigabytes int, regionSlug string) (*godo.Action, *godo.Response, error)
}

func (f *StorageActions) Attach(volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
	if f.AttachFn == nil {
		return nil, nil, nil
	}
	return f.AttachFn(volumeID, dropletID)
}

func (f *StorageActions) Detach(volumeID string) (*godo.Action, *godo.Response, error) {
	if f.DetachFn == nil {
		return nil, nil, nil
	}
	return f.DetachFn(volumeID)
}

func (f *StorageActions) Get(volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
	if f.GetFn == nil {
		return nil, nil, nil
	}
	return f.GetFn(volumeID, actionID)
}

func (f *StorageActions) List(volumeID string, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	if f.ListFn == nil {
		return nil, nil, nil
	}
	return f.ListFn(volumeID, opt)
}

func (f *StorageActions) Resize(volumeID string, sizeGigabytes int, regionSlug string) (*godo.Action, *godo.Response, error) {
	if f.ResizeFn == nil {
		return nil, nil, nil
	}
	return f.ResizeFn(volumeID, sizeGigabytes, regionSlug)
}
