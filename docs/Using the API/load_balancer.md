# Load Balancer

## Example

Supergiant currently supports adding load balancers via its API as in the example below:

### Request

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

### Response

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