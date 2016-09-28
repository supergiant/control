################################################################################
# NOTE This is using cmd/cli/cli.go, built as a `supergiant` binary, which can #
#      be done using the following command:                                    #
#                                                                              #
#      go build -o $GOPATH/bin/supergiant cmd/cli/cli.go                       #
#                                                                              #
# You'll then need to globally configure server address and API token:         #
#                                                                              #
#      supergiant configure -s <server_address> -t <api_token>                 #
#                                                                              #
################################################################################

set -e

: ${KUBE_NAME?"Need to set KUBE_NAME"}
: ${ENTRYPOINT_NAME?"Need to set ENTRYPOINT_NAME"}

# Create an internal Service for clustering / transport.
cat <<EOF | supergiant kube_resources create -f -
{
  "kube_name": "$KUBE_NAME",
  "namespace": "my-es-cluster",
  "kind": "Service",
  "name": "transport",
  "template": {
    "spec": {
      "selector": {
        "service": "elasticsearch"
      },
      "ports": [
        {
          "port": 9300
        }
      ]
    }
  }
}
EOF

# Create a NodePort Service for external HTTP access. We use the special
# Supergiant EntrypointListener definition to link this Service to our
# Entrypoint on a specified port.
cat <<EOF | supergiant kube_resources create -f -
{
  "kube_name": "$KUBE_NAME",
  "namespace": "my-es-cluster",
  "kind": "Service",
  "name": "http",
  "template": {
    "spec": {
      "type": "NodePort",
      "selector": {
        "service": "elasticsearch"
      },
      "ports": [
        {
          "name": "http",
          "port": 9200,
          "SUPERGIANT_ENTRYPOINT_LISTENER": {
            "entrypoint_name": "$ENTRYPOINT_NAME",
            "entrypoint_port": 80
          }
        }
      ]
    }
  }
}
EOF

# Create a Pod for the first Elasticsearch node. We use the special Supergiant
# Volume definition to specify an external Volume to be dynamically provisioned
# and assigned to this Pod.
cat <<EOF | supergiant kube_resources create -f -
{
  "kube_name": "$KUBE_NAME",
  "namespace": "my-es-cluster",
  "kind": "Pod",
  "name": "data-node-0",
  "template": {
    "metadata": {
      "labels": {
        "service": "elasticsearch"
      }
    },
    "spec": {
      "containers": [
        {
          "name": "elasticsearch",
          "image": "elasticsearch:2.3.3",
          "command": [
            "elasticsearch",
            "-Des.insecure.allow.root=true",
            "-Des.discovery.zen.ping.unicast.hosts=transport.my-es-cluster.svc.cluster.local:9300",
            "-Des.path.data=/data-0,/data-1",
            "-Des.path.logs=/data-0"
          ],
          "resources": {
            "requests": {
              "memory": "1.5Gi"
            }
          },
          "volumeMounts": [
            {
              "name": "es-node0-disk0",
              "mountPath": "/data-0"
            },
            {
              "name": "es-node0-disk1",
              "mountPath": "/data-1"
            }
          ]
        }
      ],
      "volumes": [
        {
          "name": "es-node0-disk0",
          "SUPERGIANT_EXTERNAL_VOLUME": {
            "type": "gp2",
            "size": 10
          }
        },
        {
          "name": "es-node0-disk1",
          "SUPERGIANT_EXTERNAL_VOLUME": {
            "type": "gp2",
            "size": 10
          }
        }
      ]
    }
  }
}
EOF

# This will print out the Entrypoint information, so you can get the ES cluster address:
# supergiant entrypoints get --id=$ENTRYPOINT_ID

# Use this command to view logs:
# supergiant kubectl -k 1 logs data-node-0 --namespace=my-es-cluster
