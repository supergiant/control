package helm

import (
	"fmt"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
)

const (
	TillerPort = 44134

	TillerConnectionTimeout int64 = 300
)

var (
	tillerPodLabels = labels.Set{"app": "helm", "name": "tiller"}
)

// Interface is an extended interface for the helm client.
type Interface interface {
	helm.Interface

	Close()
}

type client struct {
	*helm.Client

	tun *tunnel
}

// New creates a new helm client.
func New(kclient kubernetes.Interface, restConf *rest.Config, tillerNamespace string) (Interface, error) {
	tillerPodName, err := getTillerPodName(kclient.CoreV1(), tillerNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "get tiller")
	}

	tun := newTunnel(kclient.CoreV1().RESTClient(), restConf, tillerNamespace, tillerPodName, TillerPort)
	if err := tun.forwardPort(); err != nil {
		return nil, errors.Wrap(err, "setup port forwarding")
	}

	hclient := helm.NewClient(helm.Host(fmt.Sprintf("127.0.0.1:%d", tun.Local)), helm.ConnectTimeout(TillerConnectionTimeout))

	return &client{hclient, tun}, nil
}

// Close disconnects an underlying tunnel connection.
func (c *client) Close() {
	if c.tun != nil {
		c.tun.close()
	}
}

func getTillerPodName(client corev1.PodsGetter, namespace string) (string, error) {
	pod, err := getFirstRunningPod(client, namespace, tillerPodLabels.AsSelector())
	if err != nil {
		return "", err
	}

	return pod.ObjectMeta.GetName(), nil
}

func getFirstRunningPod(client corev1.PodsGetter, namespace string, selector labels.Selector) (*apiv1.Pod, error) {
	pods, err := client.Pods(namespace).List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, errors.New("could not find tiller")
	}

	for _, p := range pods.Items {
		if isPodReady(&p) {
			return &p, nil
		}
	}

	return nil, errors.New("could not find a ready tiller pod")
}
