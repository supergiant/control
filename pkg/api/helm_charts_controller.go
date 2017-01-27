package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListHelmCharts(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.HelmChart), new(model.HelmChartList))
}

func CreateHelmChart(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmChart)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmCharts.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateHelmChart(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.HelmChart)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmCharts.Update(id, new(model.HelmChart), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetHelmChart(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmChart)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmCharts.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteHelmChart(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmChart)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmCharts.Delete(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
