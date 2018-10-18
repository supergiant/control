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
