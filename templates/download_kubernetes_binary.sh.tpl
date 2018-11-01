#!/bin/bash

source /etc/environment
sudo curl -sSL -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v{{ .K8SVersion }}/bin/{{ .OperatingSystem }}/{{ .Arch }}/kubectl
sudo chmod +x /usr/bin/$FILE
sudo chmod +x /usr/bin/kubectl