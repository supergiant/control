#!/bin/bash
wget http://storage.googleapis.com/kubernetes-helm/helm-v2.6.2-linux-amd64.tar.gz --directory-prefix=/tmp/
tar -C /tmp -xvf /tmp/helm-v2.6.2-linux-amd64.tar.gz
cp /tmp/linux-amd64/helm /opt/bin/helm
chmod +x /opt/bin/helm
/opt/bin/helm init