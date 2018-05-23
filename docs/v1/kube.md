# Kube

A Kube represents a Kubernetes cluster. It belongs to a [CloudAccount](cloud_account.md), and it is the parent of [Nodes](node.md), [LoadBalancers](load_balancer.md), and [KubeResources](kube_resource.md). It encompasses all hardware-related assets such as instance types used to create masters and nodes, networking, and applications deployed to the cluster. 

Configuration for Kubes differs across cloud providers but all Kubes share a definition of Node Sizes  (machines/instance types the [Capacity Service](capacity_service.md) uses to create servers and auto-scale the cluster). 

### Creating Kubes

#### Amazon AWS

![](../img/kube-deploy_aws.gif)



#### DigitalOcean

![](../img/digital-ocean-deploy.gif)

#### 

See [Deploying Clusters](https://github.com/supergiant/supergiant/wiki/Deploying-Kubes) guide for a detailed information about cluster configuration and deployment for each supported provider. 

