# Kube

A Kube represents a Kubernetes cluster. It belongs to a
[CloudAccount](cloud_account.md), and is the parent of
[Nodes](node.md), [LoadBalancers](load_balancer.md), and
[KubeResources](kube_resource.md).
In other words, it is the encompassing object for all hardware-related assets.

### Examples

_Note the node_sizes field corresponds to what server sizes the
[Capacity Service](capacity_service.md) will use when creating servers._

#### AWS

```json
{
  "cloud_account_name": "my-aws-account",
  "name": "my-aws-kube",
  "master_node_size": "m4.large",
  "node_sizes": [
    "m4.large",
    "m4.xlarge",
    "m4.2xlarge",
    "m4.4xlarge"
  ],
  "aws_config": {
    "region": "us-east-1",
    "availability_zone": "us-east-1b"
  }
}
```

#### DigitalOcean

```json
{
  "cloud_account_name": "my-do-account",
  "name": "my-do-kube",
  "master_node_size": "2gb",
  "node_sizes": [
    "2gb",
    "4gb",
    "8gb",
    "16gb"
  ],
  "digitalocean_config": {
    "region": "nyc1",
    "ssh_key_fingerprint": "<your_do_ssh_key_fingerprint>"
  }
}
```
