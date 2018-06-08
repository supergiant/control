#!/bin/bash
  source /etc/environment
  K8S_VERSION=v{{ .GCEConfig.KubernetesVersion }}
  mkdir -p /opt/bin
  FILE=$1
  if [ ! -f /opt/bin/$FILE ]; then
    curl -sSL -o /opt/bin/$FILE https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/$FILE
    curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl
    chmod +x /opt/bin/$FILE
    chmod +x /opt/bin/kubectl
  else
    # we check the version of the binary
    INSTALLED_VERSION=$(/opt/bin/$FILE --version)
    MATCH=$(echo "${INSTALLED_VERSION}" | grep -c "${K8S_VERSION}")
    if [ $MATCH -eq 0 ]; then
      # the version is different
      curl -sSL -o /opt/bin/$FILE https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/$FILE
      curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl
      chmod +x /opt/bin/$FILE
      chmod +x /opt/bin/kubectl
    fi
  fi