package core_test

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"
)

func TestCapacityServicePerform(t *testing.T) {
	Convey("CapacityService Perform works correctly", t, func() {
		table := []struct {
			// Mocks / Input
			kubes             []*model.Kube
			providerNodeSizes map[string][]*core.NodeSize

			mockNodeExistingPods map[string][]*kubernetes.Pod // map of Node Name => []Pod
			mockPendingPods      []*kubernetes.Pod
			mockPodEvents        map[string][]*kubernetes.Event // map of Pod Name => []Event

			mockListPodsError   error
			mockListEventsError error

			mockNodeCreateError error
			mockNodeDeleteError error

			mockDBFindError error

			// Expectations
			nodeSizesCreated []string
			nodeNamesDeleted []string
			err              error
		}{
			// A full-featured successful example
			{
				// Mocks
				kubes: []*model.Kube{
					{
						CloudAccount: &model.CloudAccount{
							// Provider is used to fetch NodeSizes off of Core
							Provider: "test-provider",
						},
						Name: "test-kube",
						NodeSizes: []string{
							"1gb-test-size",
							"2gb-test-size",
							"4gb-test-size",
							"8gb-test-size",
							"16gb-test-size",
						},
						// These aren't preloaded, but we can still mock them like this.
						// NOTE this will be destroyed because we are not mocking any active
						// Pods on this with reserved resources, and the timestamp will
						// default to sometime several decades ago.
						Nodes: []*model.Node{
							{
								KubeName: "test-kube",
								Name:     "node-1.biz",
								Size:     "1gb-test-size",
							},
						},
					},
				},
				providerNodeSizes: map[string][]*core.NodeSize{
					"test-provider": {
						{
							Name:     "1gb-test-size",
							RAMGIB:   1,
							CPUCores: 1,
						},
						{
							Name:     "2gb-test-size",
							RAMGIB:   2,
							CPUCores: 1,
						},
						{
							Name:     "4gb-test-size",
							RAMGIB:   4,
							CPUCores: 1,
						},
						{
							Name:     "8gb-test-size",
							RAMGIB:   8,
							CPUCores: 2,
						},
						{
							Name:     "16gb-test-size",
							RAMGIB:   16,
							CPUCores: 4,
						},
					},
				},
				mockPendingPods: []*kubernetes.Pod{
					// Pod 1
					{
						Metadata: kubernetes.Metadata{
							Name: "test-pod",
						},
						Spec: kubernetes.PodSpec{
							Containers: []kubernetes.Container{
								{
									Resources: kubernetes.Resources{
										Requests: kubernetes.ResourceValues{
											CPU:    "1",
											Memory: "4Gi",
										},
										Limits: kubernetes.ResourceValues{
											CPU:    "2",
											Memory: "8Gi",
										},
									},
								},
							},
							Volumes: []kubernetes.Volume{
								{
									Name:                 "test-ebs-volume",
									AwsElasticBlockStore: &kubernetes.AwsElasticBlockStore{},
								},
							},
						},
					},
					// Pod 2
					{
						Metadata: kubernetes.Metadata{
							Name: "test-pod-2",
						},
						Spec: kubernetes.PodSpec{
							Containers: []kubernetes.Container{
								{
									Resources: kubernetes.Resources{
										Requests: kubernetes.ResourceValues{
											CPU:    "1",
											Memory: "4Gi",
										},
									},
								},
							},
							Volumes: []kubernetes.Volume{
								{
									Name:                 "test-ebs-volume",
									AwsElasticBlockStore: &kubernetes.AwsElasticBlockStore{},
								},
							},
						},
					},
				},
				mockPodEvents: map[string][]*kubernetes.Event{
					"test-pod": {
						{
							Message: "failed to fit in any node",
						},
					},
					"test-pod-2": {
						{
							Message: "failed to fit in any node",
						},
					},
				},
				mockNodeExistingPods: map[string][]*kubernetes.Pod{
					"existing-node": {
						{
							Metadata: kubernetes.Metadata{
								Name: "test-pod",
							},
							Spec: kubernetes.PodSpec{
								Containers: []kubernetes.Container{
									{
										Resources: kubernetes.Resources{
											// No reservation here, so this Pod will be moved off
											Limits: kubernetes.ResourceValues{
												CPU:    "1",
												Memory: "0.5Gi",
											},
										},
									},
								},
							},
						},
					},
				},
				// Expectations
				nodeSizesCreated: []string{"16gb-test-size"},
				nodeNamesDeleted: []string{"node-1.biz"},
				err:              nil,
			},

			// On DB Find error
			{
				mockDBFindError: errors.New("DBFindError"),
				err:             errors.New("DBFindError"),
			},

			// On ListPods error
			{
				kubes: []*model.Kube{
					{
						CloudAccount: &model.CloudAccount{
							// Provider is used to fetch NodeSizes off of Core
							Provider: "test-provider",
						},
						Name: "test-kube",
						NodeSizes: []string{
							"1gb-test-size",
						},
					},
				},
				providerNodeSizes: map[string][]*core.NodeSize{
					"test-provider": {
						{
							Name:     "1gb-test-size",
							RAMGIB:   1,
							CPUCores: 1,
						},
					},
				},
				mockListPodsError: errors.New("ListPodsError"),
				err:               errors.New("Capacity service error when fetching incoming pods: ListPodsError"),
			},

			// On ListEvents error
			{
				kubes: []*model.Kube{
					{
						CloudAccount: &model.CloudAccount{
							// Provider is used to fetch NodeSizes off of Core
							Provider: "test-provider",
						},
						Name: "test-kube",
						NodeSizes: []string{
							"1gb-test-size",
						},
					},
				},
				providerNodeSizes: map[string][]*core.NodeSize{
					"test-provider": {
						{
							Name:     "1gb-test-size",
							RAMGIB:   1,
							CPUCores: 1,
						},
					},
				},
				mockPendingPods: []*kubernetes.Pod{
					// Pod 1
					{
						Metadata: kubernetes.Metadata{
							Name: "test-pod",
						},
					},
				},
				mockPodEvents: map[string][]*kubernetes.Event{
					"test-pod": {
						{
							Message: "failed to fit in any node",
						},
					},
				},
				mockListEventsError: errors.New("ListEventsError"),
				err:                 errors.New("Capacity service error when fetching incoming pods: ListEventsError"),
			},

			// On NodeCreate error
			{
				kubes: []*model.Kube{
					{
						CloudAccount: &model.CloudAccount{
							// Provider is used to fetch NodeSizes off of Core
							Provider: "test-provider",
						},
						Name: "test-kube",
						NodeSizes: []string{
							"1gb-test-size",
						},
					},
				},
				providerNodeSizes: map[string][]*core.NodeSize{
					"test-provider": {
						{
							Name:     "1gb-test-size",
							RAMGIB:   1,
							CPUCores: 1,
						},
					},
				},
				mockPendingPods: []*kubernetes.Pod{
					// Pod 1
					{
						Metadata: kubernetes.Metadata{
							Name: "test-pod",
						},
					},
				},
				mockPodEvents: map[string][]*kubernetes.Event{
					"test-pod": {
						{
							Message: "failed to fit in any node",
						},
					},
				},
				mockNodeCreateError: errors.New("NodeCreateError"),
				nodeSizesCreated:    []string{"1gb-test-size"},
				err:                 errors.New("Capacity service error when creating Node: NodeCreateError"),
			},

			// On NodeDelete error
			{
				kubes: []*model.Kube{
					{
						CloudAccount: &model.CloudAccount{
							// Provider is used to fetch NodeSizes off of Core
							Provider: "test-provider",
						},
						Name: "test-kube",
						Nodes: []*model.Node{
							{
								Name: "existing-node",
							},
						},
						NodeSizes: []string{
							"1gb-test-size",
						},
					},
				},
				providerNodeSizes: map[string][]*core.NodeSize{
					"test-provider": {
						{
							Name:     "1gb-test-size",
							RAMGIB:   1,
							CPUCores: 1,
						},
					},
				},
				mockNodeExistingPods: map[string][]*kubernetes.Pod{
					"existing-node": {
						{
							Metadata: kubernetes.Metadata{
								Name: "test-pod",
							},
						},
					},
				},
				mockNodeDeleteError: errors.New("NodeDeleteError"),
				nodeNamesDeleted:    []string{"existing-node"},
				err:                 errors.New("Capacity service error when deleting Node: NodeDeleteError"),
			},
		}

		for _, item := range table {

			var nodeSizesCreated []string
			var nodeNamesDeleted []string

			c := &core.Core{
				Log: logrus.New(),

				Settings: core.Settings{
					NodeSizes: item.providerNodeSizes,
				},

				// Database
				DB: &fake_core.DB{
					FindFn: func(out interface{}, where ...interface{}) error {
						switch reflect.TypeOf(out).String() {
						case "*[]*model.Kube":
							// Set out to Kubes... NOTE we don't need to do this for Nodes,
							// since we have "preloaded" them above
							reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.kubes))
						case "*[]*model.Node":
						}
						return item.mockDBFindError
					},
				},

				// Nodes
				Nodes: &fake_core.Nodes{
					CreateFn: func(m *model.Node) error {
						nodeSizesCreated = append(nodeSizesCreated, m.Size)
						return item.mockNodeCreateError
					},
					DeleteFn: func(_ *int64, m *model.Node) core.ActionInterface {
						return &fake_core.Action{
							NowFn: func() error {
								nodeNamesDeleted = append(nodeNamesDeleted, m.Name)
								return item.mockNodeDeleteError
							},
						}
					},
				},

				// Kubernetes
				K8S: func(kube *model.Kube) kubernetes.ClientInterface {
					return &fake_core.KubernetesClient{
						ListPodsFn: func(query string) ([]*kubernetes.Pod, error) {
							// for Node hasPodsWithReservedResources
							rxp := regexp.MustCompile("fieldSelector=spec.nodeName=([^,]+),status.phase=Running")
							if rxp.MatchString(query) {
								nodeName := rxp.FindStringSubmatch(query)[1]
								return item.mockNodeExistingPods[nodeName], item.mockListPodsError
							}
							// else we are listing for pending Pods
							return item.mockPendingPods, item.mockListPodsError
						},
						ListEventsFn: func(query string) ([]*kubernetes.Event, error) {
							podName := regexp.MustCompile("involvedObject.name=(.+)").FindStringSubmatch(query)[1]
							return item.mockPodEvents[podName], item.mockListEventsError
						},
					}
				},
			}

			service := &core.CapacityService{
				Core:            c,
				WaitBeforeScale: 0,
			}

			err := service.Perform()

			So(err, ShouldResemble, item.err)
			So(nodeSizesCreated, ShouldResemble, item.nodeSizesCreated)
			So(nodeNamesDeleted, ShouldResemble, item.nodeNamesDeleted)
		}
	})
}
