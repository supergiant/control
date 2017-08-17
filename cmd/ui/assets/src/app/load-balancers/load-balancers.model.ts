export class LoadBalancersModel {
  model = {
    "kube_name": "",
    "name": "",
    "namespace": "default",
    "ports": {
      "80": 8080
    },
    "selector": {
      "key": "value"
    }
  }
  schema = {
  "properties": {
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
    "ports": {
      "properties": {
        "80": {
          "description": "80",
          "type": "string"
        }
      },
      "type": "object"
    },
    "selector": {
      "properties": {
        "key": {
          "default": "value",
          "description": "Key",
          "type": "string"
        }
      },
      "type": "object"
    }
  }
}
}
