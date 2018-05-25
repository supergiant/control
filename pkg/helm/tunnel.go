package helm

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// These functions are adapted from the "helm" repository's file
//
//    kubernetes/helm/pkg/kube/tunnel.go
//
// Reason: original package depends on large amount of the dependencies

// tunnel describes a ssh-like tunnel to a kubernetes pod
type tunnel struct {
	Local     int
	Remote    int
	Namespace string
	PodName   string
	Out       io.Writer
	stopChan  chan struct{}
	readyChan chan struct{}
	config    *rest.Config
	client    rest.Interface
}

func newTunnel(client rest.Interface, config *rest.Config, namespace, podName string, remote int) *tunnel {
	return &tunnel{
		config:    config,
		client:    client,
		Namespace: namespace,
		PodName:   podName,
		Remote:    remote,
		stopChan:  make(chan struct{}, 1),
		readyChan: make(chan struct{}, 1),
		Out:       ioutil.Discard,
	}
}

// close disconnects a tunnel connection
func (t *tunnel) close() {
	close(t.stopChan)
}

// forwardPort opens a tunnel to a kubernetes pod
func (t *tunnel) forwardPort() error {
	// Build a url to the portforward endpoint
	// example: http://localhost:8080/api/v1/namespaces/helm/pods/tiller-deploy-9itlq/portforward
	u := t.client.Post().
		Resource("pods").
		Namespace(t.Namespace).
		Name(t.PodName).
		SubResource("portforward").URL()

	transport, upgrader, err := spdy.RoundTripperFor(t.config)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", u)

	local, err := getAvailablePort()
	if err != nil {
		return fmt.Errorf("could not find an available port: %s", err)
	}
	t.Local = local

	ports := []string{fmt.Sprintf("%d:%d", t.Local, t.Remote)}

	pf, err := portforward.New(dialer, ports, t.stopChan, t.readyChan, t.Out, t.Out)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- pf.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		return fmt.Errorf("forwarding ports: %v", err)
	case <-pf.Ready:
		return nil
	}
}

func getAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}
