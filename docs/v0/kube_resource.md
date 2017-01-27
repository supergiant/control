# Kube Resource

A Kube Resource represents any object that can be created within a
[Kube](kube.md).

### Examples

#### A basic internal Service

```json
{
  "kube_name": "my-kube",
  "namespace": "default",
  "kind": "Service",
  "name": "my-private-svc",
  "resource": {
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

#### A Pod

```json
{
  "kube_name": "my-kube",
  "namespace": "my-namespace",
  "kind": "Pod",
  "name": "my-pod",
  "resource": {
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
          "name": "just-a-dir",
          "emptyDir": {}
        }
      ]
    }
  }
}
```
