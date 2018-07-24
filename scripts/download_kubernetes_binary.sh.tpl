#!/bin/bash
source /etc/environment
K8S_VERSION=v{{ .KubernetesVersion }}
mkdir -p /opt/bin
curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/{{ .OperatingSystem }}/{{ .Arch }}/kubectl
chmod +x /opt/bin/$FILE
chmod +x /opt/bin/kubectl