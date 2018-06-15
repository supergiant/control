#!/bin/bash
echo "{{ .ConfigFile }}" > /etc/kubernetes/config.json
sudo docker run --privileged=true --volume=/etc/ssl/cer:/usr/share/ca-certificates --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl gcr.io/google_containers/hyperkube:v1.8.7 /hyperkube proxy --config /etc/kubernetes/config.json --master http://{{ .Kube.MasterPrivateIP }}
