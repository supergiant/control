# Load Balancer

Load Balancer is a Kube Resource that allows adding load balancers to a Kubernetes cluster. Supergiant API supports  [external CSP-specific load balancers](https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/) of a type "**LoadBalancer**" defined as a Service and  HTTP(s) load balancers defined via **Ingress** resource. 

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

Currently, Supergiant supports adding load balancers via its API like in the example below:

#### Request

```json
{
  "kube_name": "my-kube",
  "name": "my-entrypoint",
  "namespace": "default",
  "selector": {
    "app": "my-app"
  },
  "ports": {
    "80": 8080,
    "443": 8081
  }
}
```

#### Response

```json
{
  "kube_name": "my-kube",
  "name": "my-entrypoint",
  "namespace": "default",
  "selector": {
    "app": "my-app"
  },
  "ports": {
    "80": 8080,
    "443": 8081
  },
  "address": "elb.blah.blah.amazonaws.com"
}
```

