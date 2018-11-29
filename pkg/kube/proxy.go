package kube

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type APIProxy struct {
	servicesMux sync.RWMutex
	// map[kubeID]map[serviceID]*ServiceProxy
	k8sClusterIDToProxies map[string]map[string]*ServiceProxy
	svc                   Interface
	logger                logrus.FieldLogger
}

type ServiceProxy struct {
	service     *ServiceInfo
	srv         *http.Server
	servingBase string
}

func NewAPIProxy(svc Interface, logger logrus.FieldLogger) *APIProxy {
	return &APIProxy{
		k8sClusterIDToProxies: make(map[string]map[string]*ServiceProxy),
		svc:                   svc,
		logger:                logger,
	}
}

func NewServiceProxy(service *ServiceInfo, logger logrus.FieldLogger) (*ServiceProxy, error) {
	var httpServer = &http.Server{}
	var mux = http.NewServeMux()
	httpServer.Handler = mux

	url, err := url.Parse(service.APIServerProxyURL)
	if err != nil {
		return nil, err
	}
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	mux.HandleFunc("/", newHandler(service, proxy, logger))

	// web server on a randomly (os) chosen port
	listener, err := net.Listen("tcp", "")
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

	return &ServiceProxy{
		service:     service,
		srv:         httpServer,
		servingBase: addr,
	}, nil
}

func (p *APIProxy) SetServices(kubeID string, services []*ServiceInfo) error {
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

		proxy, err := NewServiceProxy(service, p.logger)
		if err != nil {
			p.logger.Errorf("can't create proxy for serviceID: %v, err: %v", err, service.ID)
			continue
		}
		p.k8sClusterIDToProxies[kubeID][service.ID] = proxy
	}
	p.servicesMux.Unlock()

	for k := range p.k8sClusterIDToProxies[kubeID] {
		p.logger.Errorf("%+v", *p.k8sClusterIDToProxies[kubeID][k])
	}

	return nil
}

func (p *APIProxy) GetServicesPorts(kubeID string) map[string]string {
	var serviceIDToPort = map[string]string{}

	p.servicesMux.RLock()

	for serviceID, serviceProxy := range p.k8sClusterIDToProxies[kubeID] {
		var parts = strings.Split(serviceProxy.servingBase, ":")
		serviceIDToPort[serviceID] = parts[len(parts)-1]
	}
	p.servicesMux.RUnlock()

	return serviceIDToPort
}

func newHandler(service *ServiceInfo, reverseProxy *httputil.ReverseProxy, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		var inputURL = req.URL.Path

		inputURL = strings.TrimPrefix(inputURL, req.URL.Scheme)
		inputURL = strings.TrimPrefix(inputURL, req.URL.Host)

		if strings.HasPrefix(inputURL, service.SelfLink) {
			inputURL = strings.TrimPrefix(inputURL, service.SelfLink)
		}

		if strings.Index(inputURL, ":") == 0 {
			parts := strings.Split(inputURL, "/")
			inputURL = strings.TrimPrefix(inputURL, parts[0]+"/proxy")
		}

		logger.Infof("req.URL: %+v, inputURL %s, inputFullURL: %s, service: %+v",
			req.URL,
			inputURL,
			req.URL.String(),
			*service,
		)

		req.SetBasicAuth(service.User, service.Password)

		// Update the headers to allow for SSL redirection
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.URL.Path = inputURL

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		reverseProxy.ServeHTTP(res, req)
	}
}

func (p *APIProxy) Shutdown(ctx context.Context) {
	p.servicesMux.Lock()
	defer p.servicesMux.Unlock()
	for clusterID := range p.k8sClusterIDToProxies {
		for serviceID, proxy := range p.k8sClusterIDToProxies[clusterID] {
			if err := proxy.Shutdown(ctx); err != nil {
				p.logger.Errorf("cant close server for serviceID: %v, error: %v", serviceID, err)
			}
		}
	}
}

func (sp *ServiceProxy) Shutdown(ctx context.Context) error {
	if err := sp.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
