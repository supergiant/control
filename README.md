```shell
# Run etcd in background (or in another tab/pane)
etcd &

ETCD_ENDPOINT=http://localhost:2379 \
AWS_REGION=us-east-1 \
AWS_AZ=us-east-1c \
AWS_SG_ID=<security_group_id> \
AWS_SUBNET_ID=<subnet_id> \
K8S_HOST=<kube_master_ip> \
K8S_USER=<kube_http_basic_username> \
K8S_PASS=<kube_http_basic_password> \
go run main.go
```

See [example.sh](example.sh) and [api/router.go](api/router.go).
