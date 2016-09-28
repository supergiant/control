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

# Create an external Service and Pod for each member
for member in {0..2}
do
# Service
cat <<EOF | supergiant kube_resources create -f -
{
  "kube_name": "$KUBE_NAME",
  "namespace": "my-mongo-db",
  "kind": "Service",
  "name": "member-$member",
  "template": {
    "spec": {
      "type": "NodePort",
      "selector": {
        "member": "$member"
      },
      "ports": [
        {
          "name": "mongo-$member",
          "port": 27017,
          "SUPERGIANT_ENTRYPOINT_LISTENER": {
            "entrypoint_name": "$ENTRYPOINT_NAME",
            "entrypoint_port": $((33333 + $member))
          }
        }
      ]
    }
  }
}
EOF

# Pod
cat <<EOF | supergiant kube_resources create -f -
{
  "kube_name": "$KUBE_NAME",
  "namespace": "my-mongo-db",
  "kind": "Pod",
  "name": "member-$member",
  "template": {
    "metadata": {
      "labels": {
        "member": "$member"
      }
    },
    "spec": {
      "containers": [
        {
          "name": "mongo",
          "image": "mongo",
          "command": [
            "mongod", "--replSet", "rs0"
          ],
          "resources": {
            "requests": {
              "memory": "0.25Gi"
            },
            "limits": {
              "cpu": 0.25,
              "memory": "1Gi"
            }
          },
          "volumeMounts": [
            {
              "name": "mongodata-$member",
              "mountPath": "/data/db"
            }
          ]
        }
      ],
      "volumes": [
        {
          "name": "mongodata-$member",
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
done

echo "Waiting for Pod to start"
while [[ $(supergiant kube_resources list --filter=kind:Pod --filter=name:member-0 --format='{{ .Started }}') == 'false' ]]; do
  printf .
  sleep 1
done

# Get external address of first member
first_member_address=$(supergiant entrypoints list --filter=name:$ENTRYPOINT_NAME --format='{{ .Address }}:33334')

# Configure Replica Set
mongo $first_member_address --eval 'rs.initiate(); rs.reconfig({
  _id: "rs0",
  members: [
    {_id: 0, host: "member-0.my-mongo-db.svc.cluster.local:27017"},
    {_id: 1, host: "member-1.my-mongo-db.svc.cluster.local:27017"},
    {_id: 2, host: "member-2.my-mongo-db.svc.cluster.local:27017"}
  ]
})'

echo "You can reach Mongo at $first_member_address"

kube_id=$(supergiant kubes list --filter=name:$KUBE_NAME --format='{{ .ID }}')
echo "And you can view the container log with:"
echo "supergiant kubectl -k $kube_id logs data-node-0 --namespace=my-es-cluster"
