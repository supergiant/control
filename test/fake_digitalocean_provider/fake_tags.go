package fake_digitalocean_provider

import "github.com/digitalocean/godo"

type Tags struct {
	ListFn           func(*godo.ListOptions) ([]godo.Tag, *godo.Response, error)
	GetFn            func(string) (*godo.Tag, *godo.Response, error)
	CreateFn         func(*godo.TagCreateRequest) (*godo.Tag, *godo.Response, error)
	UpdateFn         func(string, *godo.TagUpdateRequest) (*godo.Response, error)
	DeleteFn         func(string) (*godo.Response, error)
	TagResourcesFn   func(string, *godo.TagResourcesRequest) (*godo.Response, error)
	UntagResourcesFn func(string, *godo.UntagResourcesRequest) (*godo.Response, error)
}

func (f *Tags) List(l *godo.ListOptions) ([]godo.Tag, *godo.Response, error) {
	if f.ListFn == nil {
		return nil, nil, nil
	}
	return f.ListFn(l)
}

func (f *Tags) Get(s string) (*godo.Tag, *godo.Response, error) {
	if f.GetFn == nil {
		return nil, nil, nil
	}
	return f.GetFn(s)
}

func (f *Tags) Create(r *godo.TagCreateRequest) (*godo.Tag, *godo.Response, error) {
	if f.CreateFn == nil {
		return nil, nil, nil
	}
	return f.CreateFn(r)
}

func (f *Tags) Update(s string, r *godo.TagUpdateRequest) (*godo.Response, error) {
	if f.UpdateFn == nil {
		return nil, nil
	}
	return f.UpdateFn(s, r)
}

func (f *Tags) Delete(s string) (*godo.Response, error) {
	if f.DeleteFn == nil {
		return nil, nil
	}
	return f.DeleteFn(s)
}

func (f *Tags) TagResources(s string, r *godo.TagResourcesRequest) (*godo.Response, error) {
	if f.TagResourcesFn == nil {
		return nil, nil
	}
	return f.TagResourcesFn(s, r)
}

func (f *Tags) UntagResources(s string, r *godo.UntagResourcesRequest) (*godo.Response, error) {
	if f.UntagResourcesFn == nil {
		return nil, nil
	}
	return f.UntagResourcesFn(s, r)
}
