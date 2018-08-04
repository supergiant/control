# Supergiant

<img src="http://supergiant.io/img/logo_dark.svg" width="400">

---

Supergiant empowers developers and administrators through its simplified deployment and management of Kubernetes, in addition to easing the configuration and deployment of Helm charts, taking advantage of Kubernetes' power, flexibility, and abstraction.

Supergiant facilitates clusters on multiple cloud providers, striving for truly agnostic and impartial infrastructure--and it does this with an autoscaling system that cares deeply about efficiency. It asserts through downscaling and resource packing that unutilized infrastructure shouldn't be paid for (and, therefore, shouldn't be running).

Supergiant implements simple practices that abstract load-balancing, application deployment, basic monitoring, node deployment or destruction, and more, on a highly usable UI. Its efficient packing algorithm enables seamless auto-scaling of Kubernetes clusters, minimizing costs while maintaining the resiliency of applications. To dive into top-level concepts, see [the documentation](https://supergiant.readthedocs.io/en/v1.0.0/API/capacity_service/).

# Features

* Fully compatible with native Kubernetes versions 1.5.7, 1.6.7, 1.7.7, and 1.8.7
* Easy management and deployment of multiple kubes in various configurations
* AWS, DigitalOcean, OpenStack, Packet, GCE, and on-premise kube deployment
* Easy creation of Helm releases, Pods, Services, LoadBalancers, etc.
* Automatic, resource-based node scaling
* Compatibility with multiple hardware architectures
* Role-based Users, Session-based logins, self-signed SSLs, and API tokens
* A clean UI and CLI, both built on top of an API (with importable [Go client lib](pkg/client))