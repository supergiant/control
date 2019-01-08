package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type ReverseProxyContainer struct {
	logger           logrus.FieldLogger
	ProxiesPortRange PortRange

	servicesMux sync.RWMutex
	// map[kubeID + serviceID]*ServiceReverseProxy
	Proxies     map[string]*ServiceReverseProxy
	PortToProxy map[int32]*ServiceReverseProxy
}

type PortRange struct {
	From, To int32
}

type Target struct {
	ProxyID    string
	TargetURL  string
	KubeConfig *rest.Config
}

type Container interface {
	RegisterProxies(targets []Target) error
	GetProxies(prefix string) map[string]*ServiceReverseProxy
	Shutdown(ctx context.Context)
}

type ServiceReverseProxy struct {
	TargetURL   string
	ServingBase string

	srv *http.Server
}

func NewReverseProxyContainer(proxiesPortRange PortRange, logger logrus.FieldLogger) *ReverseProxyContainer {
	return &ReverseProxyContainer{
		Proxies:          make(map[string]*ServiceReverseProxy),
		logger:           logger,
		ProxiesPortRange: proxiesPortRange,
	}
}

func NewServiceProxy(port int32, targetURL string, tr http.RoundTripper, logger logrus.FieldLogger) (*ServiceReverseProxy, error) {
	var httpServer = &http.Server{}
	var mux = http.NewServeMux()
	httpServer.Handler = mux

	url, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = tr

	mux.HandleFunc("/", proxyHandler(url.Path, proxy, logger))

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
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

	logger.Infof("proxy server started on: %s, for targetURL: %+v", addr, targetURL)

	return &ServiceReverseProxy{
		ServingBase: addr,

		srv: httpServer,
	}, nil
}

func (p *ReverseProxyContainer) RegisterProxies(targets []Target) error {
	if p.Proxies == nil {
		return errors.New("component was not initialized properly")
	}

	p.servicesMux.Lock()
	defer p.servicesMux.Unlock()

	for _, t := range targets {
		_, exists := p.Proxies[t.ProxyID]
		if exists {
			continue
		}

		prx, err := p.register(t)
		if err != nil {
			p.logger.Errorf("can't create proxy for serviceID: %v, err: %v", err, t.ProxyID)
			continue
		}

		p.Proxies[t.ProxyID] = prx
	}

	return nil
}

func (p *ReverseProxyContainer) GetProxies(prefix string) map[string]*ServiceReverseProxy {
	var result = map[string]*ServiceReverseProxy{}
	p.servicesMux.RLock()
	defer p.servicesMux.RUnlock()

	for proxyID, proxy := range p.Proxies {
		if !strings.HasPrefix(proxyID, prefix) {
			continue
		}
		result[proxyID] = proxy
	}

	return result
}

func (p *ReverseProxyContainer) register(t Target) (*ServiceReverseProxy, error) {
	if t.KubeConfig == nil {
		return nil, errors.New("rest config should be provided")
	}

	port, err := getNonOccupiedPort(p.ProxiesPortRange)
	if err != nil {
		return nil, err
	}

	tr, err := rest.TransportFor(t.KubeConfig)
	if err != nil {
		return nil, err
	}

	return NewServiceProxy(port, t.TargetURL, tr, p.logger)
}

func proxyHandler(baseuri string, reverseProxy *httputil.ReverseProxy, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Debugf("baseuri: %s, req.URL: %+v, inputURL %s",
			baseuri,
			req.URL.Path,
			strings.TrimPrefix(req.URL.Path, baseuri),
		)
		// Path prefix has been set to proxy
		req.URL.Path = strings.TrimPrefix(req.URL.Path, baseuri)

		// Update the headers to allow for SSL redirection
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		reverseProxy.ServeHTTP(res, req)
	}
}

func (p *ReverseProxyContainer) Shutdown(ctx context.Context) {
	p.servicesMux.Lock()
	defer p.servicesMux.Unlock()

	for proxyID, proxy := range p.Proxies {
		if err := proxy.shutdown(ctx); err != nil {
			p.logger.Errorf("cant close server for proxyID: %v, error: %v", proxyID, err)
		}
	}

}

func (sp *ServiceReverseProxy) shutdown(ctx context.Context) error {
	if err := sp.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (sp *ServiceReverseProxy) Port() string {
	var parts = strings.Split(sp.ServingBase, ":")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func getNonOccupiedPort(portRange PortRange) (int32, error) {
	for i := portRange.From; i < portRange.To; i++ {
		if checkPort(i) {
			return i, nil
		}
	}

	return 0, errors.Errorf("there is no non occupied port in range from %d to %d", portRange.From, portRange.To)
}

func checkPort(port int32) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
