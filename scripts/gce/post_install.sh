#!/bin/bash
until $(curl --output /dev/null --silent --head --fail http://127.0.0.1:8080); do printf '.'; sleep 5; done
curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://127.0.0.1:8080/api/v1/namespaces
/opt/bin/kubectl config set-cluster default-cluster --server="127.0.0.1:8080"
/opt/bin/kubectl config set-context default-system --cluster=default-cluster --user=default-admin
/opt/bin/kubectl config use-context default-system
/opt/bin/kubectl create -f /etc/kubernetes/addons/kube-dns.yaml
/opt/bin/kubectl create -f /etc/kubernetes/addons/cluster-monitoring
/opt/bin/helm init