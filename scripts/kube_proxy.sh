#!/bin/bash
sudo docker run --privileged=true \
--volume=/etc/ssl/certs:/usr/share/ca-certificates \
--volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro \
--volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl \
gcr.io/google_containers/hyperkube:v1.8.7 /hyperkube proxy \
--master=https://{{ .Kube.MasterPrivateIP }} \
--kubeconfig=/etc/kubernetes/worker-kubeconfig.yaml \
--v=2 --proxy-mode=iptables