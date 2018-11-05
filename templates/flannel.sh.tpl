#!/bin/bash

sudo wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
sudo mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld
sudo chmod 755 /usr/bin/flanneld

sudo bash -c "cat << EOF > /etc/systemd/system/flanneld.service
[Unit]
Description=Networking service

[Service]
Restart=always

Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
Environment="ETCDCTL_API=3"
ExecStart=/usr/bin/flanneld --etcd-endpoints=http://{{ .EtcdHost }}:2379

[Install]
WantedBy=multi-user.target
EOF"
sudo systemctl daemon-reload
sudo systemctl enable flanneld.service
sudo systemctl start flanneld.service

while [ ! -f /run/flannel/subnet.env ]; do   sleep 2; printf '.'; done
sudo source /run/flannel/subnet.env

sudo cat << EOF > /etc/systemd/system/docker.service
[Unit]
Requires=flanneld.service
After=flanneld.service

[Service]
Restart=always
ExecStart=/usr/bin/dockerd  --bip=${FLANNEL_SUBNET} --mtu=${FLANNEL_MTU}
EOF

sudo systemctl stop docker.service
sudo systemctl daemon-reload
sudo systemctl restart docker.service