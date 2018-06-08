#!/bin/bash
source /etc/environment
K8S_VERSION=v{{ .Kube.KubernetesVersion }}
mkdir -p /opt/bin
mkdir /etc/multipath/
touch /etc/multipath/bindings
curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl
chmod +x /opt/bin/$FILE
chmod +x /opt/bin/kubectl

curl -sSL -o /opt/bin/cni.tar.gz https://storage.googleapis.com/kubernetes-release/network-plugins/cni-07a8a28637e97b22eb8dfe710eeae1344f69d16e.tar.gz
tar xzf "/opt/bin/cni.tar.gz" -C "/opt/bin" --overwrite
mv /opt/bin/bin/* /opt/bin
rm -r /opt/bin/bin/
rm -f "/opt/bin/cni.tar.gz"

cd /opt/bin/
git clone https://github.com/packethost/packet-block-storage.git
cd packet-block-storage
chmod 755 ./*
/opt/bin/packet-block-storage/packet-block-storage-attach

cd /tmp
wget https://github.com/digitalocean/doctl/releases/download/v1.4.0/doctl-1.4.0-linux-amd64.tar.gz
tar xf /tmp/doctl-1.4.0-linux-amd64.tar.gz
sudo mv /tmp/doctl /opt/bin/
sudo mkdir -p /root/.config/doctl/
sudo touch /root/.config/doctl/config.yaml