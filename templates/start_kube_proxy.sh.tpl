#!/bin/bash
echo '{
        "bind-address": "{{ .MasterPrivateIP }}",
        "hostname-override": "{{ .MasterPrivateIP }}",
        "cluster-cidr": "172.30.0.0/16",
        "logtostderr": true,
        "v": 0,
        "allow-privileged": true,
        "master": "http://{{ .MasterPrivateIP }}:{{ .ProxyPort }}",
        "etcd-servers": "http://{{ .MasterPrivateIP }}:{{ .EtcdClientPort }}"
      }
' > /etc/kubernetes/config.json
sudo docker run --privileged=true --volume=/etc/ssl/cer:/usr/share/ca-certificates --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl gcr.io/google_containers/hyperkube:v{{ .K8SVersion }} /hyperkube proxy --config /etc/kubernetes/config.json --master http://{{ .MasterPrivateIP }}
