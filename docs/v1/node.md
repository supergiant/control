# Node

A Node is a pairing of a server (i.e. VM, an EC2 instance, a DigitalOcean droplet etc.) from a Provider with a Kubernetes Node object. One of the main goals of Supergiant is to abstract away server management entirely. Therefore, the [Capacity Service](capacity_service.md) is capable of managing servers autonomously so that users can focus on allocating containers. At the same time, however, users can manually add Nodes to an existing cluster whenever they like via Supergiant UI or CRUD API.

### Adding Nodes Via Supergiant UI

To create a new node, first select the cluster you want to add the node to. Then, in the "**Nodes**" section of the cluster resources, select a node type (e.g `m4.large` for Amazon EC2) from the dropdown list and click "**Create Node**". A new node will be added to your cluster in a trice. 

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/concepts/create-new-node.gif)

### Adding Nodes Via Supergiant API

```json
{
  "kube_name": "my-kube",
  "size": "c4.large"
}
```
