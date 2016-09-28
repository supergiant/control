# Entrypoint

An Entrypoint is an external load balancer, that can be linked to containers
through the use of [EntrypointListeners](entrypoint_listener.md).

While Kubernetes does provide LoadBalancer Services, we found them to be less
flexible than desired -- we wanted non-namespaced load balancers that could
vary in implementation as we support more clouds.

### Examples

#### Request

```json
{
  "kube_name": "my-kube",
  "name": "my-entrypoint"
}
```

#### Response

```json
{
  "kube_name": "my-kube",
  "name": "my-entrypoint",
  "provider_id": "839403",
  "address": "elb.blah.blah.amazonaws.com"
}
```
