#!/bin/bash
curl -sSL -o /opt/bin/cni.tar.gz https://storage.googleapis.com/kubernetes-release/network-plugins/cni-07a8a28637e97b22eb8dfe710eeae1344f69d16e.tar.gz
tar xzf "/opt/bin/cni.tar.gz" -C "/opt/bin" --overwrite
mv /opt/bin/bin/* /opt/bin
rm -r /opt/bin/bin/
rm -f "/opt/bin/cni.tar.gz"