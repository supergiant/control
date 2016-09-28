# Entrypoint Listener

An Entrypoint Listener is a port mapping to route public internet traffic to a
container. `entrypoint_port` represents the external visitor-facing port, and
`node_port` represents the target the external port routes to.

### Examples

#### Direct API usage

```json
{
  "entrypoint_name": "my-entrypoint",
  "entrypoint_port": 80,
  "entrypoint_protocol": "HTTP",
  "node_port": 30333,
  "node_protocol": "TCP"
}
```

#### In a KubeResource Template for a NodePort Service

_Note that in this usage, the node_port can be automatically assigned, which is
the recommended usage pattern. The port value is used to map the node_port to
the container (representing the port your container uses), and is not used by
the EntrypointListener._

```json
{
  "kube_name": "my-kube",
  "namespace": "default",
  "kind": "Service",
  "name": "my-service",
  "template": {
    "spec": {
      "type": "NodePort",
      "selector": {
        "service": "my-service"
      },
      "ports": [
        {
          "name": "my-port",
          "port": 8080,
          "SUPERGIANT_ENTRYPOINT_LISTENER": {
            "entrypoint_name": "my-entrypoint",
            "entrypoint_protocol": "HTTP",
            "entrypoint_port": 80
          }
        }
      ]
    }
  }
}
```
