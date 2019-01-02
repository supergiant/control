package proxy

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func TestNewReverseProxyContainer(t *testing.T) {
	portRange := PortRange{
		0,
		65535,
	}
	logger := logrus.New()
	reverseProxy := NewReverseProxyContainer(portRange, logger)

	if reverseProxy.logger != logger {
		t.Errorf("Wrong logger expected %v actual %v",
			logger, reverseProxy.logger)
	}

	if reverseProxy.ProxiesPortRange.From != portRange.From ||
		reverseProxy.ProxiesPortRange.To != portRange.To {
		t.Errorf("Wrong values of portRange %v expected %v",
			reverseProxy.ProxiesPortRange, portRange)
	}
}

func TestServiceReverseProxy_Port(t *testing.T) {
	expectedPort := "8080"
	p := &ServiceReverseProxy{
		ServingBase: fmt.Sprintf("host:%s", expectedPort),
	}

	actualPort := p.Port()

	if expectedPort != actualPort {
		t.Errorf("Wrong port expected %s actual %s",
			expectedPort, actualPort)
	}
}

func TestServiceReverseProxy_PortErr(t *testing.T) {
	p := &ServiceReverseProxy{
		ServingBase: "",
	}

	actualPort := p.Port()

	if actualPort != "" {
		t.Errorf("Wrong port expected empty "+
			"string actual %s", actualPort)
	}
}

func TestReverseProxyContainer_GetProxies(t *testing.T) {
	reverseProxy := &ReverseProxyContainer{
		Proxies: map[string]*ServiceReverseProxy{
			"key1":        {},
			"key2":        {},
			"prefixHello": {},
		},
	}

	proxies := reverseProxy.GetProxies("prefix")

	if len(proxies) != 1 {
		t.Errorf("Wrong len of proxies "+
			"expected %d actual %d", 1, len(proxies))
	}
}

func TestReverseProxyContainer_Shutdown(t *testing.T) {
	proxy := &ReverseProxyContainer{
		Proxies: map[string]*ServiceReverseProxy{
			"proxy1": {
				srv: &http.Server{},
			},
			"proxy2": {
				srv: &http.Server{},
			},
		},
	}

	ctx := context.Background()
	proxy.Shutdown(ctx)
}

func TestReverseProxyContainer_ShutdownErr(t *testing.T) {
	proxy := &ReverseProxyContainer{
		Proxies: map[string]*ServiceReverseProxy{
			"proxy1": {
				srv: &http.Server{},
			},
		},
	}

	deadline := time.Now().Add(-time.Second)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	proxy.Shutdown(ctx)
}

func TestNewHandler(t *testing.T) {
	u := &url.URL{}
	reverseproxy := httputil.NewSingleHostReverseProxy(u)
	logger := logrus.New()
	h := newHandler("selflink",
		"user", "password", reverseproxy, logger)

	if h == nil {
		t.Errorf("Handler value must not be nil")
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet,
		"http://hostname.com/selflink", nil)

	h(rec, req)
}

func TestCheckPort(t *testing.T) {
	isOccupied := checkPort(-1)

	if isOccupied {
		t.Errorf("Wrong value port expected "+
			"false actual %v", isOccupied)
	}
}

func TestGetNonOccupiedPortErr(t *testing.T) {
	portRange := PortRange{
		-1,
		-100,
	}

	_, err := getNonOccupiedPort(portRange)

	if err == nil {
		t.Errorf("error must not be nil")
	}
}

func TestGetNonOccupiedPortSuccess(t *testing.T) {
	portRange := PortRange{
		1024,
		65535,
	}

	port, err := getNonOccupiedPort(portRange)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if port > portRange.To || port < portRange.From {
		t.Errorf("Port %d is out of port range [%d,%d]",
			port, portRange.From, portRange.To)
	}
}

func TestNewServiceProxy(t *testing.T) {
	logger := logrus.New()
	proxy, err := NewServiceProxy(8080, "/url",
		"selflink", "user", "password", logger)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if proxy == nil {
		t.Errorf("proxy must not be nil")
	}
}

func TestReverseProxyContainer_RegisterProxiesErr(t *testing.T) {
	containter := &ReverseProxyContainer{}
	err := containter.RegisterProxies([]*Target{})

	if err == nil {
		t.Errorf("error must be nil")
	}
}

func TestReverseProxyContainer_RegisterProxies(t *testing.T) {
	containter := &ReverseProxyContainer{
		logger: logrus.New(),
		Proxies: map[string]*ServiceReverseProxy{
			"1234": {},
		},
		ProxiesPortRange: PortRange{
			1024,
			65535,
		},
	}
	err := containter.RegisterProxies([]*Target{
		{
			ProxyID: "1234",
		},
		{
			ProxyID:   "5678",
			TargetURL: "/url",
			SelfLink:  "selflink",
			User:      "user",
			Password:  "password",
		},
	})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
