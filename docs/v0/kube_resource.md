# Kube Resource

A Kube Resource represents any object that can be created with a
[Kube](kube.md). It extends the Kubernetes API slightly by providing
two special types of template objects that instruct Supergiant to allocate and
assign external cloud assets
([EntrypointListeners](entrypoint_listener.md) for NodePort Services and
[Volumes](volumes.md) for Pods).

We employ this system of templating and preprocessing to provide more
flexibility in our implementation of persistent storage and external load
balancing for containers, which is a big need as we continue to support more
cloud providers.

Kube Resources can be **started** and **stopped**, which allows for changes and
manual application through restart (which gives the user deployment flow
flexibility). Stopping will not delete any assets (if the KubeResource owns
any), it will just delete the KubeResource within the Kube. That way, you don't
lose your persistent volume or load balancer port allocation during a restart.

### Examples

#### A basic internal Service

```json
{
  "kube_name": "my-kube",
  "namespace": "default",
  "kind": "Service",
  "name": "my-private-svc",
  "template": {
    "spec": {
      "selector": {
        "service": "my-pod-selector"
      },
      "ports": [
        {
          "port": 8080
        }
      ]
    }
  }
}
```

#### A NodePort Service with EntrypointListener definition

```json
{
  "kube_name": "my-kube",
  "namespace": "default",
  "kind": "Service",
  "name": "my-public-svc",
  "template": {
    "spec": {
      "type": "NodePort",
      "selector": {
        "service": "my-pod-selector"
      },
      "ports": [
        {
          "name": "website-http",
          "port": 8080,
          "SUPERGIANT_ENTRYPOINT_LISTENER": {
            "entrypoint_name": "my-entrypoint",
            "entrypoint_port": 80
          }
        }
      ]
    }
  }
}
```

#### A Pod with external Volume definition

```json
{
  "kube_name": "my-kube",
  "namespace": "my-namespace",
  "kind": "Pod",
  "name": "my-pod",
  "template": {
    "metadata": {
      "labels": {
        "service": "my-pod-selector",
      }
    },
    "spec": {
      "containers": [
        {
          "name": "my-container",
          "image": "some/image:v0.1.0",
          "volumeMounts": [
            {
              "name": "disk-0",
              "mountPath": "/mnt"
            },
            {
              "name": "just-a-dir",
              "mountPath": "/some_dir"
            }
          ]
        }
      ],
      "volumes": [
        {
          "name": "disk-0",
          "SUPERGIANT_EXTERNAL_VOLUME": {
            "type": "gp2",
            "size": 20
          }
        },
        {
          "name": "just-a-dir",
          "emptyDir": {}
        }
      ]
    }
  }
}
```
