package kubernetes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKubernetesEnsureNamespace(t *testing.T) {
	Convey("Kubernetes EnsureNamespace works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			namespace string
			// Mocks
			mockGetNamespaceResponseCode int
			mockGetNamespaceResponseBody string

			mockCreateNamespaceResponseCode int
			mockCreateNamespaceResponseBody string

			// Expectations
			namespaceNameCreated string
			err                  error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				namespace: "test",
				// Mocks
				mockGetNamespaceResponseCode:    404,
				mockCreateNamespaceResponseCode: 201,
				// Expectations
				namespaceNameCreated: "test",
				err:                  nil,
			},

			// When Namespace already exists
			{
				// Input
				kube:      &model.Kube{},
				namespace: "test",
				// Mocks
				mockGetNamespaceResponseCode: 200,
				// Expectations
				err: nil,
			},

			// On unexpected Get error
			{
				// Input
				kube:      &model.Kube{},
				namespace: "test",
				// Mocks
				mockGetNamespaceResponseCode: 500,
				mockGetNamespaceResponseBody: "crud",
				// Expectations
				err: errors.New("K8S 500 error: crud"),
			},

			// On unexpected Create error
			{
				// Input
				kube:      &model.Kube{},
				namespace: "test",
				// Mocks
				mockGetNamespaceResponseCode:    404,
				mockCreateNamespaceResponseCode: 666,
				mockCreateNamespaceResponseBody: "the devil",
				// Expectations
				namespaceNameCreated: "test",
				err:                  errors.New("K8S 666 error: the devil"),
			},
		}

		for _, item := range table {

			var namespaceNameCreated string

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/namespaces/"+item.namespace {
								// GetNamespace
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockGetNamespaceResponseCode),
									StatusCode: item.mockGetNamespaceResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockGetNamespaceResponseBody)),
								}

							} else if r.Method == "POST" && r.URL.Path == "/api/v1/namespaces" {

								reqBody, _ := ioutil.ReadAll(r.Body)
								defer r.Body.Close()
								namespace := &kubernetes.Namespace{}
								json.Unmarshal(reqBody, namespace)

								namespaceNameCreated = namespace.Metadata.Name

								// CreateNamespace
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockCreateNamespaceResponseCode),
									StatusCode: item.mockCreateNamespaceResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockCreateNamespaceResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path)
							}
							return

						},
					},
				},
			}

			err := kubernetes.EnsureNamespace(item.namespace)

			So(err, ShouldResemble, item.err)
			So(namespaceNameCreated, ShouldEqual, item.namespaceNameCreated)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesGetResource(t *testing.T) {
	Convey("Kubernetes GetResource works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			kind      string
			namespace string
			name      string
			// Mocks
			mockGetResourceResponseCode int
			mockGetResourceResponseBody string
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Pod",
				namespace: "test",
				name:      "testname",
				// Mocks
				mockGetResourceResponseCode: 200,
				mockGetResourceResponseBody: `{}`,
				// Expectations
				err: nil,
			},

			// On error
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Pod",
				namespace: "test",
				name:      "testname",
				// Mocks
				mockGetResourceResponseCode: 404,
				mockGetResourceResponseBody: `unexpected error`,
				// Expectations
				err: errors.New("K8S 404 error: unexpected error"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/namespaces/"+item.namespace+"/"+strings.ToLower(item.kind)+"s/"+item.name {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockGetResourceResponseCode),
									StatusCode: item.mockGetResourceResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockGetResourceResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path)
							}
							return

						},
					},
				},
			}

			out := json.RawMessage([]byte(`{}`))
			err := kubernetes.GetResource("api/v1", item.kind, item.namespace, item.name, &out)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesCreateResource(t *testing.T) {
	Convey("Kubernetes CreateResource works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			kind      string
			namespace string
			in        map[string]interface{}
			// Mocks
			mockCreateResourceResponseCode int
			mockCreateResourceResponseBody string
			// Expectations
			resourceNameCreated string
			err                 error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Service",
				namespace: "test",
				in: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-resource",
					},
					"spec": map[string]interface{}{
						"type": "NodePort",
					},
				},
				// Mocks
				mockCreateResourceResponseCode: 201,
				mockCreateResourceResponseBody: `{}`,
				// Expectations
				resourceNameCreated: "my-resource",
				err:                 nil,
			},

			// On 409 error (already exists)
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Service",
				namespace: "test",
				in: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-resource",
					},
					"spec": map[string]interface{}{
						"type": "NodePort",
					},
				},
				// Mocks
				mockCreateResourceResponseCode: 409,
				mockCreateResourceResponseBody: `{}`,
				// Expectations
				resourceNameCreated: "my-resource",
				err:                 nil,
			},

			// On other error
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Service",
				namespace: "test",
				in: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-resource",
					},
					"spec": map[string]interface{}{
						"type": "NodePort",
					},
				},
				// Mocks
				mockCreateResourceResponseCode: 500,
				mockCreateResourceResponseBody: `bad thing`,
				// Expectations
				resourceNameCreated: "my-resource",
				err:                 errors.New("K8S 500 error: bad thing"),
			},
		}

		for _, item := range table {

			var resourceNameCreated string

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "POST" && r.URL.Path == "/api/v1/namespaces/"+item.namespace+"/"+strings.ToLower(item.kind)+"s" {

								reqBody, _ := ioutil.ReadAll(r.Body)
								defer r.Body.Close()
								resource := &kubernetes.Namespace{} // NOTE we're just getting name, so it doesn't matter what type we use here
								json.Unmarshal(reqBody, resource)

								resourceNameCreated = resource.Metadata.Name

								resp = &http.Response{
									Status:     strconv.Itoa(item.mockCreateResourceResponseCode),
									StatusCode: item.mockCreateResourceResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockCreateResourceResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path)
							}
							return

						},
					},
				},
			}

			out := json.RawMessage([]byte(`{}`))
			err := kubernetes.CreateResource("api/v1", item.kind, item.namespace, item.in, &out)

			So(err, ShouldResemble, item.err)
			So(resourceNameCreated, ShouldEqual, item.resourceNameCreated)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesDeleteResource(t *testing.T) {
	Convey("Kubernetes DeleteResource works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			kind      string
			namespace string
			name      string
			// Mocks
			mockDeleteResourceResponseCode int
			mockDeleteResourceResponseBody string
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Pod",
				namespace: "test",
				name:      "testname",
				// Mocks
				mockDeleteResourceResponseCode: 202,
				mockDeleteResourceResponseBody: `{}`,
				// Expectations
				err: nil,
			},

			// On error
			{
				// Input
				kube:      &model.Kube{},
				kind:      "Pod",
				namespace: "test",
				name:      "testname",
				// Mocks
				mockDeleteResourceResponseCode: 404,
				mockDeleteResourceResponseBody: `not found`,
				// Expectations
				err: errors.New("K8S 404 error: not found"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "DELETE" && r.URL.Path == "/api/v1/namespaces/"+item.namespace+"/"+strings.ToLower(item.kind)+"s/"+item.name {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockDeleteResourceResponseCode),
									StatusCode: item.mockDeleteResourceResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockDeleteResourceResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path)
							}
							return

						},
					},
				},
			}

			err := kubernetes.DeleteResource("api/v1", item.kind, item.namespace, item.name)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesListNamespaces(t *testing.T) {
	Convey("Kubernetes ListNamespaces works correctly", t, func() {
		table := []struct {
			// Input
			kube  *model.Kube
			query string
			// Mocks
			mockListNamespacesResponseCode int
			mockListNamespacesResponseBody string
			// Expectations
			items []*kubernetes.Namespace
			err   error
		}{
			// A successful example
			{
				// Input
				kube:  &model.Kube{},
				query: "some=query",
				// Mocks
				mockListNamespacesResponseCode: 200,
				mockListNamespacesResponseBody: `{
          "items": [
            {
              "metadata": {
                "name": "namespacer"
              }
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.Namespace{
					{
						Metadata: kubernetes.Metadata{
							Name: "namespacer",
						},
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListNamespacesResponseCode: 500,
				mockListNamespacesResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/namespaces" && r.URL.Query().Encode() == item.query {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListNamespacesResponseCode),
									StatusCode: item.mockListNamespacesResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListNamespacesResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListNamespaces(item.query)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesListNodes(t *testing.T) {
	Convey("Kubernetes ListNodes works correctly", t, func() {
		table := []struct {
			// Input
			kube  *model.Kube
			query string
			// Mocks
			mockListNodesResponseCode int
			mockListNodesResponseBody string
			// Expectations
			items []*kubernetes.Node
			err   error
		}{
			// A successful example
			{
				// Input
				kube:  &model.Kube{},
				query: "some=query",
				// Mocks
				mockListNodesResponseCode: 200,
				mockListNodesResponseBody: `{
          "items": [
            {
              "metadata": {
                "name": "Node"
              }
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.Node{
					{
						Metadata: kubernetes.Metadata{
							Name: "Node",
						},
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListNodesResponseCode: 500,
				mockListNodesResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/nodes" && r.URL.Query().Encode() == item.query {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListNodesResponseCode),
									StatusCode: item.mockListNodesResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListNodesResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListNodes(item.query)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesListPods(t *testing.T) {
	Convey("Kubernetes ListPods works correctly", t, func() {
		table := []struct {
			// Input
			kube  *model.Kube
			query string
			// Mocks
			mockListPodsResponseCode int
			mockListPodsResponseBody string
			// Expectations
			items []*kubernetes.Pod
			err   error
		}{
			// A successful example
			{
				// Input
				kube:  &model.Kube{},
				query: "some=query",
				// Mocks
				mockListPodsResponseCode: 200,
				mockListPodsResponseBody: `{
          "items": [
            {
              "metadata": {
                "name": "Pod"
              }
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.Pod{
					{
						Metadata: kubernetes.Metadata{
							Name: "Pod",
						},
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListPodsResponseCode: 500,
				mockListPodsResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/pods" && r.URL.Query().Encode() == item.query {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListPodsResponseCode),
									StatusCode: item.mockListPodsResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListPodsResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListPods(item.query)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesListEvents(t *testing.T) {
	Convey("Kubernetes ListEvents works correctly", t, func() {
		table := []struct {
			// Input
			kube  *model.Kube
			query string
			// Mocks
			mockListEventsResponseCode int
			mockListEventsResponseBody string
			// Expectations
			items []*kubernetes.Event
			err   error
		}{
			// A successful example
			{
				// Input
				kube:  &model.Kube{},
				query: "some=query",
				// Mocks
				mockListEventsResponseCode: 200,
				mockListEventsResponseBody: `{
          "items": [
            {
              "metadata": {
                "name": "Event"
              }
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.Event{
					{
						Metadata: kubernetes.Metadata{
							Name: "Event",
						},
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListEventsResponseCode: 500,
				mockListEventsResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/events" && r.URL.Query().Encode() == item.query {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListEventsResponseCode),
									StatusCode: item.mockListEventsResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListEventsResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListEvents(item.query)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------

func TestKubernetesListPodHeapsterCPUUsageMetrics(t *testing.T) {
	Convey("Kubernetes ListPodHeapsterCPUUsageMetrics works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			namespace string
			name      string
			// Mocks
			mockListPodHeapsterCPUUsageMetricsResponseCode int
			mockListPodHeapsterCPUUsageMetricsResponseBody string
			// Expectations
			items []*kubernetes.HeapsterMetric
			err   error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				namespace: "namespace",
				name:      "name",
				// Mocks
				mockListPodHeapsterCPUUsageMetricsResponseCode: 200,
				mockListPodHeapsterCPUUsageMetricsResponseBody: `{
          "metrics": [
            {
              "value": 9999
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.HeapsterMetric{
					{
						Value: 9999,
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListPodHeapsterCPUUsageMetricsResponseCode: 500,
				mockListPodHeapsterCPUUsageMetricsResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+item.namespace+"/pods/"+item.name+"/metrics/cpu-usage" {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListPodHeapsterCPUUsageMetricsResponseCode),
									StatusCode: item.mockListPodHeapsterCPUUsageMetricsResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListPodHeapsterCPUUsageMetricsResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListPodHeapsterCPUUsageMetrics(item.namespace, item.name)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesListPodHeapsterRAMUsageMetrics(t *testing.T) {
	Convey("Kubernetes ListPodHeapsterRAMUsageMetrics works correctly", t, func() {
		table := []struct {
			// Input
			kube      *model.Kube
			namespace string
			name      string
			// Mocks
			mockListPodHeapsterRAMUsageMetricsResponseCode int
			mockListPodHeapsterRAMUsageMetricsResponseBody string
			// Expectations
			items []*kubernetes.HeapsterMetric
			err   error
		}{
			// A successful example
			{
				// Input
				kube:      &model.Kube{},
				namespace: "namespace",
				name:      "name",
				// Mocks
				mockListPodHeapsterRAMUsageMetricsResponseCode: 200,
				mockListPodHeapsterRAMUsageMetricsResponseBody: `{
          "metrics": [
            {
              "value": 9999
            }
          ]
        }`,
				// Expectations
				items: []*kubernetes.HeapsterMetric{
					{
						Value: 9999,
					},
				},
				err: nil,
			},

			// On error
			{
				// Input
				kube: &model.Kube{},
				// Mocks
				mockListPodHeapsterRAMUsageMetricsResponseCode: 500,
				mockListPodHeapsterRAMUsageMetricsResponseBody: `something bad`,
				// Expectations
				err: errors.New("K8S 500 error: something bad"),
			},
		}

		for _, item := range table {

			kubernetes := &kubernetes.Client{
				Kube: item.kube,
				HTTPClient: &http.Client{
					Transport: &fake_http.RoundTripper{
						RoundTripFn: func(r *http.Request) (resp *http.Response, err error) {

							if r.Method == "GET" && r.URL.Path == "/api/v1/proxy/namespaces/kube-system/services/heapster/api/v1/model/namespaces/"+item.namespace+"/pods/"+item.name+"/metrics/memory-usage" {
								resp = &http.Response{
									Status:     strconv.Itoa(item.mockListPodHeapsterRAMUsageMetricsResponseCode),
									StatusCode: item.mockListPodHeapsterRAMUsageMetricsResponseCode,
									Body:       ioutil.NopCloser(bytes.NewBufferString(item.mockListPodHeapsterRAMUsageMetricsResponseBody)),
								}
							} else {
								panic("Did not recognize request Method / URL Path: " + r.Method + " " + r.URL.Path + " " + r.URL.Query().Encode())
							}
							return

						},
					},
				},
			}

			items, err := kubernetes.ListPodHeapsterRAMUsageMetrics(item.namespace, item.name)

			So(err, ShouldResemble, item.err)
			So(items, ShouldResemble, item.items)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesCoresFromCPUString(t *testing.T) {
	Convey("Kubernetes CoresFromCPUString works correctly", t, func() {
		table := []struct {
			// Input
			str string
			// Expectations
			cores float64
			err   error
		}{
			{
				str:   `2`,
				cores: 2,
				err:   nil,
			},
			{
				str:   `"2"`,
				cores: 2,
				err:   nil,
			},
			{
				str:   ``,
				cores: 0,
				err:   nil,
			},
			{
				str:   `""`,
				cores: 0,
				err:   errors.New(`Could not parse cores value from ""`),
			},
			{
				str:   `1500m`,
				cores: 1.5,
				err:   nil,
			},
			{
				str:   `"1500m"`,
				cores: 1.5,
				err:   nil,
			},
			{
				str:   `1500M`,
				cores: 0,
				err:   errors.New("Could not parse cores value from 1500M"),
			},
		}

		for _, item := range table {

			cores, err := kubernetes.CoresFromCPUString(item.str)

			So(err, ShouldResemble, item.err)
			So(cores, ShouldResemble, item.cores)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubernetesGiBFromMemString(t *testing.T) {
	Convey("Kubernetes GiBFromMemString works correctly", t, func() {
		table := []struct {
			// Input
			str string
			// Expectations
			gib float64
			err error
		}{
			{
				str: `2Gi`,
				gib: 2,
				err: nil,
			},
			{
				str: `"2Gi"`,
				gib: 2,
				err: nil,
			},
			{
				str: `2048Mi`,
				gib: 2,
				err: nil,
			},
			{
				str: `2097152Ki`,
				gib: 2,
				err: nil,
			},
			{
				str: `2147483648`,
				gib: 2,
				err: nil,
			},
			{
				str: ``,
				gib: 0,
				err: nil,
			},
			{
				str: `butt`,
				gib: 0,
				err: errors.New(`Bytes value butt does not match regex ^"?([0-9]+(?:\.[0-9]+)?)([KMG]i)?"?$`),
			},
		}

		for _, item := range table {

			gib, err := kubernetes.GiBFromMemString(item.str)

			So(err, ShouldResemble, item.err)
			So(gib, ShouldResemble, item.gib)
		}
	})
}
