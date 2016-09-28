# Node

A Node is a pairing of a server (i.e. VM, EC2 instance, DigitalOcean Droplet,
etc.) from a Provider, with a Kubernetes Node object.

One of the main goals of Supergiant is to abstract away server management
entirely -- while there is a full CRUD API for Nodes (_meaning you can easily
add Nodes to an existing cluster whenever you'd like_), the
[Capacity Service](capacity_service.md) is capable of managing servers
autonomously, so a user can focus on allocating containers.

### Example

```json
{
  "kube_name": "my-kube",
  "size": "c4.large"
}
```
