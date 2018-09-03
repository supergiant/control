#!/bin/bash
source /etc/environment
curl -sSL -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v{{ .K8SVersion }}/bin/{{ .OperatingSystem }}/{{ .Arch }}/kubectl
chmod +x /usr/bin/$FILE
chmod +x /usr/bin/kubectl