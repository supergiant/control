#!/bin/bash
source /etc/environment
mkdir -p /opt/bin
curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v{{ .K8SVersion }}/bin/{{ .OperatingSystem }}/{{ .Arch }}/kubectl
chmod +x /opt/bin/$FILE
chmod +x /opt/bin/kubectl