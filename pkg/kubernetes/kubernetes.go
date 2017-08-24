package kubernetes

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/supergiant/supergiant/pkg/model"
)

type ClientInterface interface {
	// EnsureNamespace creates a Kubernetes Namespace unless it already exists.
	EnsureNamespace(name string) error

	GetResource(apiVersion, kind, namespace, name string, out interface{}) error
	CreateResource(apiVersion, kind, namespace string, objIn interface{}, out interface{}) error
	UpdateResource(apiVersion, kind, namespace, name string, objIn interface{}, out interface{}) error
	DeleteResource(apiVersion, kind, namespace, name string) error

	ListNamespaces(query string) ([]*Namespace, error)
	ListEvents(query string) ([]*Event, error)
	ListNodes(query string) ([]*Node, error)
	ListPods(query string) ([]*Pod, error)
	ListServices(query string) ([]*Service, error)
	ListPersistentVolumes(query string) ([]*PersistentVolume, error)

	GetPodLog(namespace, name string) (string, error)

	ListNodeHeapsterStats(node string) ([]string, error)
	ListPodHeapsterCPUUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error)
	ListPodHeapsterRAMUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error)
	GetNodeHeapsterStats(node string, metricPath string) (HeapsterMetrics, error)
	ListKubeHeapsterStats() ([]string, error)
	GetKubeHeapsterStats(metricPath string) (HeapsterMetrics, error)
}

//------------------------------------------------------------------------------

var DefaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type Client struct {
	Kube       *model.Kube
	HTTPClient *http.Client
}

// EnsureNamespace implements the ClientInterface.
func (k *Client) EnsureNamespace(name string) error {
	// If we get a 404 here, we need to create
	err := k.requestInto("GET", "api/v1", "namespaces/"+name, nil, nil)
	if err == nil {
		return nil // already exists
	} else if !strings.Contains(err.Error(), "404") {
		return err // unexpected error
	}
	namespace := &Namespace{
		Metadata: Metadata{
			Name: name,
		},
	}
	return k.requestInto("POST", "api/v1", "namespaces", namespace, nil)
}

func (k *Client) GetResource(apiVersion, kind, namespace, name string, out interface{}) error {
	path := fmt.Sprintf("namespaces/%s/%s/%s", namespace, lowerPlural(kind), name)
	return k.requestInto("GET", apiVersion, path, nil, out)
}

func (k *Client) CreateResource(apiVersion, kind, namespace string, in interface{}, out interface{}) error {
	path := fmt.Sprintf("namespaces/%s/%s", namespace, lowerPlural(kind))
	err := k.requestInto("POST", apiVersion, path, in, out)
	// Only return error if it's NOT a 409 already exists error
	if err != nil && !strings.Contains(err.Error(), "409") {
		return err
	}
	return nil
}

func (k *Client) UpdateResource(apiVersion, kind, namespace, name string, in interface{}, out interface{}) error {
	path := fmt.Sprintf("namespaces/%s/%s/%s", namespace, lowerPlural(kind), name)
	return k.patchRequestInto(apiVersion, path, in, out)
}

func (k *Client) DeleteResource(apiVersion, kind, namespace, name string) error {
	path := fmt.Sprintf("namespaces/%s/%s/%s", namespace, lowerPlural(kind), name)
	return k.requestInto("DELETE", apiVersion, path, nil, nil)
}

func (k *Client) ListNamespaces(query string) ([]*Namespace, error) {
	list := new(NamespaceList)
	if err := k.requestInto("GET", "api/v1", "namespaces?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListNodes(query string) ([]*Node, error) {
	list := new(NodeList)
	if err := k.requestInto("GET", "api/v1", "nodes?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListPods(query string) ([]*Pod, error) {
	list := new(PodList)
	if err := k.requestInto("GET", "api/v1", "pods?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListServices(query string) ([]*Service, error) {
	list := new(ServiceList)
	if err := k.requestInto("GET", "api/v1", "services?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListPersistentVolumes(query string) ([]*PersistentVolume, error) {
	list := new(PersistentVolumeList)
	if err := k.requestInto("GET", "api/v1", "persistentvolumes?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListEvents(query string) ([]*Event, error) {
	list := new(EventList)
	if err := k.requestInto("GET", "api/v1", "events?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) GetPodLog(namespace, name string) (string, error) {
	path := fmt.Sprintf("namespaces/%s/pods/%s/log", namespace, name)
	resp, err := k.request("application/json", "GET", "api/v1", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (k *Client) ListKubeHeapsterStats() ([]string, error) {
	var metrics []string
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/metrics/", nil, &metrics)
	return metrics, err
}

func (k *Client) GetKubeHeapsterStats(metricPath string) (HeapsterMetrics, error) {
	metrics := HeapsterMetrics{}
	metrics.MetricName = strings.Replace(metricPath, "/", "_", -1)
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/metrics/"+metricPath+"", nil, &metrics)
	return metrics, err
}

func (k *Client) ListNodeHeapsterStats(node string) ([]string, error) {
	var metrics []string
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/nodes/"+node+"/metrics/", nil, &metrics)
	return metrics, err
}

func (k *Client) GetNodeHeapsterStats(node string, metricPath string) (HeapsterMetrics, error) {
	metrics := HeapsterMetrics{}
	metrics.MetricName = strings.Replace(metricPath, "/", "_", -1)
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/nodes/"+node+"/metrics/"+metricPath+"", nil, &metrics)
	return metrics, err
}

func (k *Client) ListPodHeapsterRAMUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error) {
	metrics := HeapsterMetrics{}
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+namespace+"/pods/"+name+"/metrics/memory-usage", nil, &metrics)
	return metrics.Metrics, err
}

func (k *Client) ListPodHeapsterCPUUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error) {
	metrics := HeapsterMetrics{}
	err := k.requestInto("GET", "api/v1", "proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+namespace+"/pods/"+name+"/metrics/cpu-usage", nil, &metrics)
	return metrics.Metrics, err
}

// Private

func (k *Client) request(contentType, method, apiVersion, path string, in interface{}) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/%s/%s", k.Kube.MasterPublicIP, apiVersion, path)

	var body []byte
	if in != nil {
		jsonIn, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		body = jsonIn
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(k.Kube.Username, k.Kube.Password)

	req.Header.Set("Content-type", contentType)

	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Status[:2] != "20" {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("K8S %s error: %s", resp.Status, string(respBody))
	}

	return resp, nil
}

// TODO we could make this much nicer if we made request a buildable object
func (k *Client) requestIntoWithContentType(contentType, method, apiVersion, path string, in interface{}, out interface{}) error {
	resp, err := k.request(contentType, method, apiVersion, path, in)
	if err != nil {
		return err
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
}

func (k *Client) requestInto(method, apiVersion, path string, in interface{}, out interface{}) error {
	return k.requestIntoWithContentType("application/json", method, apiVersion, path, in, out)
}

func (k *Client) patchRequestInto(apiVersion, path string, in interface{}, out interface{}) error {
	return k.requestIntoWithContentType("application/merge-patch+json", "PATCH", apiVersion, path, in, out)
}

// Misc ------------------------------------------------------------------------

func CoresFromCPUString(str string) (float64, error) {
	if str == "" {
		return 0, nil
	}

	rxpMillicores := regexp.MustCompile(`^"?([0-9]+)m"?$`)      // 1000m
	rxpCores := regexp.MustCompile(`^"?([0-9]+(\.[0-9]+)?)"?$`) // 1 (can have quotes)

	getNumMatch := func(rxp *regexp.Regexp) (float64, error) {
		numberStr := rxp.FindStringSubmatch(str)[1]
		return strconv.ParseFloat(numberStr, 64)
	}

	if rxpMillicores.MatchString(str) {
		num, _ := getNumMatch(rxpMillicores)
		return num / 1000.0, nil
	}

	if rxpCores.MatchString(str) {
		num, _ := getNumMatch(rxpCores)
		return num, nil
	}

	return 0, fmt.Errorf("Could not parse cores value from %s", str)
}

func GiBFromMemString(memStr string) (float64, error) {
	if memStr == "" {
		return 0, nil
	}

	rxp := regexp.MustCompile(`^"?([0-9]+(?:\.[0-9]+)?)([KMG]i)?"?$`)

	if !rxp.MatchString(memStr) {
		return 0, fmt.Errorf(`Bytes value %s does not match regex ^"?([0-9]+(?:\.[0-9]+)?)([KMG]i)?"?$`, memStr)
	}

	match := rxp.FindStringSubmatch(memStr)

	float, _ := strconv.ParseFloat(match[1], 64)

	switch match[2] {
	case "":
		return float / 1073741824, nil
	case "Ki":
		return float / 1048576, nil
	case "Mi":
		return float / 1024, nil
	case "Gi":
		return float, nil
	}

	return 0, fmt.Errorf("Cannot parse RAM GiB from %s", memStr)
}

//------------------------------------------------------------------------------

func lowerPlural(str string) string {
	return strings.ToLower(str) + "s"
}
