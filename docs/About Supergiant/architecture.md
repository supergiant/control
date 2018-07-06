# Architecture

Supergiant is a Kubernetes-as-a-Service platform that simplifies deployment and management of Kubernetes clusters, objects, and resources. Supergiant abstracts away the Kubernetes API via an easy-to-use UI that allows linking cloud accounts, deploying apps, exposing services and monitoring Kubernetes clusters. This functionality makes it easy for DevOps teams to deploy Kubernetes without diving deep into its configuration and tuning. At the same time, however, Supergiant still provides access to low-level features of the Kubernetes API and command line tools.

Supergiant augments the Kubernetes API with a powerful Capacity Service and cost-effective packing algorithm that constitute the core of its architecture. The goal of these components is to enable easy deployment and autoscaling of Kubernetes clusters to minimize computing costs and reduce deployment times.

## Capacity Service

Servers used to provision computing resources for Kubernetes clusters may come in different flavors: cloud virtual machines, on-premise bare-metal servers, cloud-based bare-metal servers, or on-premise virtual machines. In addition, each cloud service provider (e.g., AWS, GCE, etc.) offers its users tens of server types with different cpu, memory, networking, architecture, and more. Manual selection, instantiation, and management of all these diverse server types can be a complicated task. The **Capacity Service** is thus Supergiant's way of abstracting servers from application management and deployment, allowing the user to focus on the CPU, RAM, and storage applications need —- and not on which flavor of server to use for each application.

Without the Capacity Service, running Kubernetes would force DevOps teams to manage two levels of capacity: container resource allocation and server allocation. Doing the latter manually in the dynamic container environment would be quite difficult with ever-changing resource requirements and workloads. Thus, manually managing these two levels of complexity would be akin to if AWS, for example, required its users not only to request new instances but also to ask for additional capacity in the region ahead of time! Of course, this would run against the on-demand cloud model. So, what AWS does for its instances, Supergiant's Capacity Service does for Kubernetes containers. With Supergiant's Capacity Service, Kubernetes resources can be provisioned on-demand without worrying about server capacity. Supergiant will create new nodes when apps are over capacity, and it will gently remove existing nodes when sufficiently under capacity.

## The Packing Algorithm

The Capacity Service is tightly coupled with Supergiant's **Packing Algorithm**, which decides how much capacity to allocate to a Kubernetes cluster and when to allocate it.

Supergiant's Packing Algorithm augments the default Kubernetes resource allocation strategy, introducing an efficient method for minimizing computing costs. To understand how the packing algorithm works, let's first review the Kubernetes resource allocation model.

In a nutshell, Kubernetes allocates resources to containers/pods/applications based on requests and limits specified during an object's deployment. For example, each container in a pod can specify a minimum and maximum value for compute resources (CPU and RAM) which it needs to run on a single node.

- `spec.containers[].resources.limits.cpu`
- `spec.containers[].resources.limits.memory`
- `spec.containers[].resources.requests.cpu`
- `spec.containers[].resources.requests.memory`

The kube-scheduler will use this min/max ratio to efficiently allocate resources to pods. Let's assume a scenario when two app instances with 4 CPU min / 8 CPU max ratio are running on a Kubernetes node with 8 CPUs of capacity. If one of the app instances begins to consume resources equivalent to 6 CPUs, Node 1 will be out of capacity, so the Kubernetes scheduler will attempt to move the instance to another node with sufficient capacity. That's how standard resource allocation works in Kubernetes.

Supergiant augments this model with cost-effective horizontal autoscaling that considers cloud computing costs in the equation. In general, horizontal auto-scaling is the automatic addition/removal of servers to meet the application's demand for computing resources. It's different from vertical autoscaling, in which compute resources are added/removed on a single server/node.

Going back to the case above, the Kubernetes scheduler would move the **App Instance 2** to any node with sufficient resources to run it and schedule smaller apps ("**Other App**") to Node 1.

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/concepts/packing_3.png)

However, as seen in the image above, **App Instance 2** would then live on a node with significant underutilized resources (e.g., 4 CPUs). In the real world, this would lead to paying for a 8 CPU instance, although **App Instance 2** needs only 4 CPUs. If a cluster were to have hundreds of nodes with such unused capacity, cloud infrastructure costs would grow far beyond expectations.

Supergiant's packing algorithm is designed to solve this problem. If **App Instance 2** consumes more resources than it's allowed to, Supergiant will select the node type with the exact amount of resources needed to meet the demand.

![](https://s3-ap-southeast-2.amazonaws.com/sg-github-wiki-images/concepts/packing_4.png)

Instead of allocating an 8 CPU node to **App Instance 2**, Supergiant provisions **Node 2** with 4 CPUs, which exactly matches the resource requirements of the application. In this way, Supergiant's packing algorithm ensures that applications are tightly "packed" on hardware, dramatically reducing the margin of underutilized resources.

But the benefits of Supergiant's auto-scaling go further. Supergiant can pack application instances based on the minimum request for resources they need. (**Note**: In contrast, Kubernetes would not allow an app instance to run on a server with less capacity than the app's maximum threshold.) Formally speaking, Supergiant's augmentation of the Kubernetes min/max ratio translates to the following:

- The minimum resource value becomes a maximum number of components that can fit on a node.
- The maximum resource value translates to smart decisions on when the component should be throttled or moved to another node with more resources.

These rules enable apps to occupy servers that actually don't have enough resources to fit them and to help efficiently move them to other nodes if apps start exceeding their limits. 

As mentioned above, Supergiant's packing algorithm is tightly coupled with the Capacity Service. Supergiant stores the list of available CSP instance types and **Node-Sizes** configuration specified during the cluster deployment (e.g., see the example list for DigitalOcean machine types below).

```json
"digitalocean": [
      {"name": "s-1vcpu-1gb", "ram_gib": 1, "cpu_cores": 1},
      {"name": "s-1vcpu-2gb", "ram_gib": 2, "cpu_cores": 1},
      {"name": "s-1vcpu-3gb", "ram_gib": 3, "cpu_cores": 1},
      {"name": "s-2vcpu-2gb", "ram_gib": 2, "cpu_cores": 2},
      {"name": "s-3vcpu-1gb", "ram_gib": 1, "cpu_cores": 3},
      {"name": "s-2vcpu-4gb", "ram_gib": 4, "cpu_cores": 2},
      {"name": "s-4vcpu-8gb", "ram_gib": 8, "cpu_cores": 4},
      {"name": "s-6vcpu-16gb", "ram_gib": 16, "cpu_cores": 6},
      {"name": "s-8vcpu-32gb", "ram_gib": 32, "cpu_cores": 8},
      {"name": "s-12vcpu-48gb", "ram_gib": 48, "cpu_cores": 12},
      {"name": "s-16vcpu-64gb", "ram_gib": 64, "cpu_cores": 16},
      {"name": "s-20vcpu-96gb", "ram_gib": 96, "cpu_cores": 20},
      {"name": "s-24vcpu-128gb", "ram_gib": 128, "cpu_cores": 24},
      {"name": "s-32vcpu-192gb", "ram_gib": 192, "cpu_cores": 32}
    ]
```

If the app deployed to the node consumes resources beyond its minimum level so that the total amount of resource requests exceeds the total capacity of the node, Supergiant will select an appropriate machine type that will close the supply gap and will move the app instance to it. In this way, what server flavor is used and on what node an app lives becomes unimportant. Any given cluster will always be tightly packed, so that it allocates the exact amount of resources needed while dramatically reducing cloud costs.
