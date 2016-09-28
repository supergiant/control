# Capacity Service

The Capacity Service is Supergiant's way of abstracting servers from
application management and deployments. It allows the user to focus only on how
much CPU, RAM, and disk _containers_ need, not which flavor of server to use for
each application.

If you imagine _containerization_ to be similar to _hardware-level
virtualization_ (hypervisors partitioning big host machines in the cloud into
multiple virtual machines), then the following analogy could be made:

Without the Supergiant Capacity Service, running a container orchestration
platform like Kubernetes forces devops engineers to *manage 2 levels of capacity*
-- that is, they must not only worry about container resource allocation, but
also _server allocation_. It would be equally difficult to imagine AWS requiring
its users not only to request new servers, but also request additional capacity
in the region ahead of time!

Of course, that would be silly -- cloud servers are generally _on-demand_.
That's what the Capacity Service does for Kubernetes containers. Containers can
be provisioned on-demand without worrying about server capacity. Supergiant will
handle creating Nodes when over capacity, and (gently) deleting Nodes when
sufficiently under capacity.
