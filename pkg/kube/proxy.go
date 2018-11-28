package kube

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)


type APIProxy struct {
	servicesMux      sync.RWMutex
	// map[kubeID]map[serviceID]*ServiceProxy
	k8sClusterIDToProxies map[string]map[string]*ServiceProxy
	svc              Interface
	logger logrus.FieldLogger
}

type ServiceProxy struct {
	service *ServiceInfo
	srv *http.Server
	servingBase string
}

func NewAPIProxy(svc Interface, logger logrus.FieldLogger) *APIProxy {
	return &APIProxy{
		k8sClusterIDToProxies: make(map[string]map[string]*ServiceProxy),
		svc:              svc,
		logger:logger,
	}
}

func NewServiceProxy(service *ServiceInfo, 	logger logrus.FieldLogger) (*ServiceProxy, error){
	var httpServer = &http.Server{}
	var mux = http.NewServeMux()
	httpServer.Handler = mux

	mux.HandleFunc("/", newHandler(service, logger))

	// web server on a randomly (os) chosen port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	var addr = listener.Addr().String()
	var addrParts = strings.Split(addr, ":")
	if len(addrParts) == 0 {
		return nil, errors.Errorf("can't get non occupied port, addr %v", addr)
	}

	go func() {
		if err := httpServer.Serve(listener); err != nil {
			if err != nil {
				logger.Errorf("error while serving address: %v", listener.Addr())
			}
		}
	}()

	logger.Infof("proxy server started on: %s, for service: %+v", addr, service)

	var proxy =  &ServiceProxy{
		service: service,
		srv:     httpServer,
		servingBase: addr,
	}

	return proxy, nil
}


func (p *APIProxy) SetServices(kubeID string, services []ServiceInfo) error {
	if p.svc == nil || p.k8sClusterIDToProxies == nil {
		return errors.New("component was not initialized properly")
	}

	p.servicesMux.Lock()

	_, exists := p.k8sClusterIDToProxies[kubeID]
	if !exists {
		p.k8sClusterIDToProxies[kubeID] = make(map[string]*ServiceProxy)
	}

	for _, service := range services {
		_, exists := p.k8sClusterIDToProxies[kubeID][service.ID]
		if exists {
			continue
		}

		proxy, err := NewServiceProxy(&service, p.logger)
		if err != nil {
			p.logger.Errorf("can't create proxy for serviceID: %v, err: %v", err, service.ID)
			continue
		}
		p.k8sClusterIDToProxies[kubeID][service.ID] = proxy
	}
	p.servicesMux.Unlock()

	return nil
}


// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {

	url, err := url.Parse(target)
	if err != nil {

	}
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func newHandler(service *ServiceInfo, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		var inputURL = req.URL.String()

		inputURL = strings.TrimPrefix(inputURL, req.URL.Scheme)
		inputURL = strings.TrimPrefix(inputURL, req.URL.Host)

		if strings.Contains(inputURL, service.SelfLink){
			inputURL = strings.TrimPrefix(inputURL, service.SelfLink)
		}

		var inputFullURL = req.URL.String()

		var result = service.APIServerProxyURL + inputURL


		logger.Infof("req.URL: %+v, inputURL %s, inputFullURL %s, result: %s, service: %+v",
			req.URL,
			inputURL,
			inputFullURL,
			result,
			*service,
		)

		serveReverseProxy(result, res, req)
	}
}

func (p *APIProxy)Shutdown(ctx context.Context){
	for clusterID, _ := range p.k8sClusterIDToProxies {
		for serviceID, proxy := range p.k8sClusterIDToProxies[clusterID] {
			if err := proxy.Shutdown(ctx); err != nil {
				p.logger.Errorf("cant close server for serviceID: %v, error: %v", serviceID, err)
			}
		}
	}
}

func (sp *ServiceProxy)Shutdown(ctx context.Context) error {
	if err := sp.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}