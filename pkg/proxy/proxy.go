package proxy

import (
	"context"
	"crypto/tls"
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

type Target struct{ ProxyID, TargetURL, SelfLink, User, Password string }

type Container interface {
	RegisterProxies(targets []*Target) error
	GetProxies(prefix string) map[string]*ServiceReverseProxy
	Shutdown(ctx context.Context)
}

type ServiceReverseProxy struct {
	TargetURL   string
	SelfLink    string
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

func NewServiceProxy(port int32, targetURL, selfLink, user, password string, logger logrus.FieldLogger) (*ServiceReverseProxy, error) {
	var httpServer = &http.Server{}
	var mux = http.NewServeMux()
	httpServer.Handler = mux

	url, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	mux.HandleFunc("/", newHandler(selfLink, user, password, proxy, logger))

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

func (p *ReverseProxyContainer) RegisterProxies(targets []*Target) error {
	if p.Proxies == nil {
		return errors.New("component was not initialized properly")
	}

	p.servicesMux.Lock()

	for _, targetSvc := range targets {
		_, exists := p.Proxies[targetSvc.ProxyID]
		if exists {
			continue
		}

		var port, err = getNonOccupiedPort(p.ProxiesPortRange)
		if err != nil {
			errors.Wrap(err, "can't create proxy")
		}

		p.logger.Infof("returned port: %d", port)

		proxy, err := NewServiceProxy(port, targetSvc.TargetURL, targetSvc.SelfLink, targetSvc.User, targetSvc.Password, p.logger)
		if err != nil {
			p.logger.Errorf("can't create proxy for serviceID: %v, err: %v", err, targetSvc.ProxyID)
			continue
		}
		p.Proxies[targetSvc.ProxyID] = proxy
	}
	p.servicesMux.Unlock()

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

func newHandler(selfLink, user, password string, reverseProxy *httputil.ReverseProxy, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		var inputURL = req.URL.Path

		inputURL = strings.TrimPrefix(inputURL, req.URL.Scheme)
		inputURL = strings.TrimPrefix(inputURL, req.URL.Host)

		if strings.HasPrefix(inputURL, selfLink) {
			inputURL = strings.TrimPrefix(inputURL, selfLink)
		}

		if strings.Index(inputURL, ":") == 0 {
			parts := strings.Split(inputURL, "/")
			inputURL = strings.TrimPrefix(inputURL, parts[0]+"/proxy")
		}

		logger.Infof("req.URL: %+v, inputURL %s, inputFullURL: %s",
			req.URL,
			inputURL,
			req.URL.String(),
		)

		req.SetBasicAuth(user, password)

		// Update the headers to allow for SSL redirection
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.URL.Path = inputURL

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
