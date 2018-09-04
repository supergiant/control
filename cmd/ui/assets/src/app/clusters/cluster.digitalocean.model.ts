export class ClusterDigitalOceanModel {
  digitalocean = {
    'data': {
      "clusterName": "",
      "cloudAccountName": "",
      "profile": {
          "masterProfiles": [{
              "image": "ubuntu-16-04-x64",
              "size": "s-1vcpu-2gb"
          }],
          "nodesProfiles": [{
                "image": "ubuntu-16-04-x64",
                "size": "s-2vcpu-4gb"
            },
            {
                "image": "ubuntu-16-04-x64",
                "size": "s-2vcpu-4gb"
            }
          ],
          "provider": "digitalocean",
          "region": "fra1",
          "arch": "amd64",
          "operatingSystem": "linux",
          "ubuntuVersion": "xenial",
          "dockerVersion": "17.06.0",
          "K8SVersion": "1.11.1",
          "flannelVersion": "0.10.0",
          "networkType": "vxlan",
          "cidr": "10.0.0.0/24",
          "helmVersion": "2.8.0",
          "rbacEnabled": false
      }
    },
    'schema': {
      'properties': {
        "clusterName": {
          "type": "string",
          "description": "Name of new cluster"
        },
        "cloudAccountName": {
          "type": "string",
          "description": "Select a cloud account"
        },
        "profile": {
          "type": "object",
          "properties": {
            "masterProfiles": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "image": {
                    "type": "string"
                  },
                  "size": {
                    "type": "string"
                  }
                }
              }
            },
            "nodesProfiles": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "image": {
                    "type": "string"
                  },
                  "size": {
                    "type": "string"
                  }
                }
              }
            },
            "provider": {
              "type": "string"
            },
            "region": {
              "type": "string"
            },
            "arch": {
              "type": "string"
            },
            "operatingSystem": {
              "type": "string"
            },
            "ubuntuVersion": {
              "type": "string"
            },
            "dockerVersion": {
              "type": "string"
            },
            "K8SVersion": {
              "type": "string"
            },
            "flannelVersion": {
              "type": "string"
            },
            "networkType": {
              "type": "string"
            },
            "cidr": {
              "type": "string"
            },
            "helmVersion": {
              "type": "string"
            },
            "rbacEnabled": {
              "type": "boolean"
            },
          }
        }
      }
    }
  };
}
