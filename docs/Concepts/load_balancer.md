# Load Balancer

The Load Balancer is a Kube Resource that allows adding load balancers to a Kubernetes cluster. The Supergiant API supports  [external CSP-specific load balancers](https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/) of type "**LoadBalancer**," which is defined as a Service and HTTP(S) load balancer defined via an **Ingress** resource.

1.  External CSP-specific load balancers can be added by specifying `type: LoadBalancer` on a Service manifest. See the example from the Kubernetes documentation: 

```json
    {
      "kind": "Service",
      "apiVersion": "v1",
      "metadata": {
        "name": "example-service"
      },
      "spec": {
        "ports": [{
          "port": 8765,
          "targetPort": 9376
        }],
        "selector": {
          "app": "example"
        },
        "type": "LoadBalancer"
      }
    }
```

2. Supergiant supports adding HTTP(s) load balancers using an Ingress resource. This resource can create HTTP(s) traffic rules and make context-aware load balancing decisions.