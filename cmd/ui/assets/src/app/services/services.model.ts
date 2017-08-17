export class ServicesModel {
  model = {
  "kind": "Service",
  "kube_name": "",
  "name": "",
  "namespace": "default",
  "template": {
    "spec": {
      "ports": [
        {
          "name": "jenkins",
          "port": 8080
        }
      ],
      "selector": {},
      "type": "NodePort"
    }
  }
}
schema = {
  "properties": {
    "kind": {
      "default": "Service",
      "description": "Kind",
      "type": "string"
    },
    "kube_name": {
      "description": "Kube Name",
      "type": "string"
    },
    "name": {
      "description": "Name",
      "type": "string"
    },
    "namespace": {
      "default": "default",
      "description": "Namespace",
      "type": "string"
    },
    "template": {
      "properties": {
        "spec": {
          "properties": {
            "ports": {
              "description": "Ports",
              "type": "string"
            },
            "selector": {
              "type": "string",
              "description": "Selector"
            },
            "type": {
              "default": "NodePort",
              "description": "Type",
              "type": "string"
            }
          },
          "type": "object"
        }
      },
      "type": "object"
    }
  }
}

}
