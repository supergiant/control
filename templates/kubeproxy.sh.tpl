#!/bin/bash
mkdir -p  /etc/kubernetes
sudo docker run --privileged=true --volume=/etc/ssl/certs:/usr/share/ca-certificates \
    --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro \
    --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl \
    --volume=/etc/kubernetes/config.json:/etc/kubernetes/config.json \
    gcr.io/google_containers/hyperkube:v{{ .K8SVersion }} /hyperkube proxy \
    --master http://{{ .MasterPrivateIP }}
