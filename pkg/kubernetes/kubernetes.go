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

	GetResource(kind string, namespace string, name string, out *json.RawMessage) error
	CreateResource(kind string, namespace string, objIn map[string]interface{}, out *json.RawMessage) error
	DeleteResource(kind string, namespace string, name string) error

	ListNamespaces(query string) ([]*Namespace, error)
	ListEvents(query string) ([]*Event, error)
	ListNodes(query string) ([]*Node, error)
	ListPods(query string) ([]*Pod, error)
	ListNodeHeapsterStats() ([]*HeapsterStats, error)
	ListPodHeapsterCPUUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error)
	ListPodHeapsterRAMUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error)
}

//------------------------------------------------------------------------------

// TODO
var globalK8SHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type Client struct {
	Kube *model.Kube
}

// EnsureNamespace implements the ClientInterface.
func (k *Client) EnsureNamespace(name string) error {
	// If we get a 404 here, we need to create
	err := k.requestInto("GET", "namespaces/"+name, nil, nil)
	if err == nil {
		return err // unexpected error
	} else if !strings.Contains(err.Error(), "404") {
		return nil // already exists
	}
	namespace := &Namespace{
		Metadata: Metadata{
			Name: name,
		},
	}
	return k.requestInto("POST", "namespaces", namespace, nil)
}

func (k *Client) GetResource(kind string, namespace string, name string, out *json.RawMessage) error {
	path := fmt.Sprintf("namespaces/%s/%s/%s", namespace, lowerPlural(kind), name)
	return k.requestInto("GET", path, nil, out)
}

func (k *Client) CreateResource(kind string, namespace string, in map[string]interface{}, out *json.RawMessage) error {
	path := fmt.Sprintf("namespaces/%s/%s", namespace, lowerPlural(kind))
	err := k.requestInto("POST", path, in, out)
	// Only return error if it's NOT a 409 already exists error
	if err != nil && !strings.Contains(err.Error(), "409") {
		return err
	}
	return nil
}

func (k *Client) DeleteResource(kind string, namespace string, name string) error {
	path := fmt.Sprintf("namespaces/%s/%s/%s", namespace, lowerPlural(kind), name)
	return k.requestInto("DELETE", path, nil, nil)
}

func (k *Client) ListNamespaces(query string) ([]*Namespace, error) {
	list := new(NamespaceList)
	if err := k.requestInto("GET", "namespaces?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListNodes(query string) ([]*Node, error) {
	list := new(NodeList)
	if err := k.requestInto("GET", "nodes?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListPods(query string) ([]*Pod, error) {
	list := new(PodList)
	if err := k.requestInto("GET", "pods?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListEvents(query string) ([]*Event, error) {
	list := new(EventList)
	if err := k.requestInto("GET", "events?"+query, nil, list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (k *Client) ListNodeHeapsterStats() ([]*HeapsterStats, error) {
	var metrics []*HeapsterStats
	err := k.requestInto("GET", "proxy/namespaces/kube-system/services/heapster/api/v1/model/nodes", nil, &metrics)
	return metrics, err
}

func (k *Client) ListPodHeapsterRAMUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error) {
	metrics := HeapsterMetrics{}
	err := k.requestInto("GET", "proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+namespace+"/pods/"+name+"/metrics/memory-usage", nil, &metrics)
	return metrics.Metrics, err
}

func (k *Client) ListPodHeapsterCPUUsageMetrics(namespace string, name string) ([]*HeapsterMetric, error) {
	metrics := HeapsterMetrics{}
	err := k.requestInto("GET", "proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+namespace+"/pods/"+name+"/metrics/cpu-usage", nil, &metrics)
	return metrics.Metrics, err
}

// Private

func (k *Client) requestInto(method string, path string, in interface{}, out interface{}) error {
	url := fmt.Sprintf("https://%s/api/v1/%s", k.Kube.MasterPublicIP, path)

	// fmt.Println("---------------- REQUESTING: ", method, url)

	var body []byte
	if in != nil {
		jsonIn, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = jsonIn

		// fmt.Println("---------------- BODY: ", string(body))

	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(k.Kube.Username, k.Kube.Password)

	resp, err := globalK8SHTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Status[:2] != "20" {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("K8S %s error: %s", resp.Status, string(respBody))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
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
		num, err := getNumMatch(rxpMillicores)
		if err != nil {
			return 0, err
		}
		return num / 1000.0, nil
	}

	if rxpCores.MatchString(str) {
		num, err := getNumMatch(rxpCores)
		if err != nil {
			return 0, err
		}
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

	float, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, err
	}

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
