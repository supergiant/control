#!/bin/bash
wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld
chmod 755 /usr/bin/flanneld

cat << EOF > /etc/systemd/system/flanneld.service
[Unit]
Description=Networking service

[Service]
Restart=always

Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
Environment="ETCDCTL_API=3"
ExecStart=/usr/bin/flanneld --etcd-endpoints=http://{{ .EtcdHost }}:2379

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable flanneld.service
systemctl start flanneld.service

while [ ! -f /run/flannel/subnet.env ]; do   sleep 2; printf '.'; done
source /run/flannel/subnet.env

cat << EOF > /etc/systemd/system/docker.service
[Unit]
Requires=flanneld.service
After=flanneld.service

[Service]
Restart=always
ExecStart=/usr/bin/dockerd  --bip=${FLANNEL_SUBNET} --mtu=${FLANNEL_MTU}
EOF

systemctl stop docker.service
systemctl daemon-reload
systemctl restart docker.service