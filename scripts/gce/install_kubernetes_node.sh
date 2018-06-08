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

  ## Install Tiller ##
  wget http://storage.googleapis.com/kubernetes-helm/helm-v2.6.2-linux-amd64.tar.gz --directory-prefix=/tmp/
  tar -C /tmp -xvf /tmp/helm-v2.6.2-linux-amd64.tar.gz
  cp /tmp/linux-amd64/helm /opt/bin/helm
  chmod +x /opt/bin/helm
  /opt/bin/helm init

  ## Install CNI

  curl -sSL -o /opt/bin/cni.tar.gz https://storage.googleapis.com/kubernetes-release/network-plugins/cni-07a8a28637e97b22eb8dfe710eeae1344f69d16e.tar.gz
  tar xzf "/opt/bin/cni.tar.gz" -C "/opt/bin" --overwrite
  mv /opt/bin/bin/* /opt/bin
  rm -r /opt/bin/bin/
  rm -f "/opt/bin/cni.tar.gz"

  openssl genrsa -out /etc/kubernetes/ssl/ca-key.pem 2048
  openssl req -x509 -new -nodes -key /etc/kubernetes/ssl/ca-key.pem -days 10000 -out /etc/kubernetes/ssl/ca.pem -subj "/CN=kube-ca"
  sed -e "s/\${MASTER_HOST}/$COREOS_PUBLIC_IPV4/" < /etc/kubernetes/ssl/openssl.cnf.template > /etc/kubernetes/ssl/openssl.cnf.public
  sed -e "s/\${PRIVATE_HOST}/$COREOS_PRIVATE_IPV4/" < /etc/kubernetes/ssl/openssl.cnf.public > /etc/kubernetes/ssl/openssl.cnf
  openssl genrsa -out /etc/kubernetes/ssl/apiserver-key.pem 2048
  openssl req -new -key /etc/kubernetes/ssl/apiserver-key.pem -out /etc/kubernetes/ssl/apiserver.csr -subj "/CN=kube-apiserver" -config /etc/kubernetes/ssl/openssl.cnf
  openssl x509 -req -in /etc/kubernetes/ssl/apiserver.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/apiserver.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
  openssl genrsa -out /etc/kubernetes/ssl/worker-key.pem 2048
  openssl req -new -key /etc/kubernetes/ssl/worker-key.pem -out /etc/kubernetes/ssl/worker.csr -subj "/CN=kube-worker"
  openssl x509 -req -in /etc/kubernetes/ssl/worker.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/worker.pem -days 365
  openssl genrsa -out /etc/kubernetes/ssl/admin-key.pem 2048
  openssl req -new -key /etc/kubernetes/ssl/admin-key.pem -out /etc/kubernetes/ssl/admin.csr -subj "/CN=kube-admin"
  openssl x509 -req -in /etc/kubernetes/ssl/admin.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/admin.pem -days 365
  chmod 600 /etc/kubernetes/ssl/*-key.pem
  chown root:root /etc/kubernetes/ssl/*-key.pem