```shell
# Run etcd in background (or in another tab/pane)
etcd &

ETCD_ENDPOINT=http://localhost:2379 \
AWS_REGION=us-east-1 \
AWS_AZ=us-east-1c \
K8S_HOST=<kube_master_ip \
K8S_USER=<kube_http_basic_username> \
K8S_PASS=<kube_http_basic_password> \
go run main.go
```

See [example.sh](example.sh) and [controller/router.go](controller/router.go).

*Note: I'm going to squash commit history soon to get rid of the private
Dockerhub key in example.sh.*
