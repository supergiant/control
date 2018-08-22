#!/bin/bash
wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld
chmod 755 /usr/bin/flanneld

# install etcdctl
GITHUB_URL=https://github.com/coreos/etcd/releases/download
ETCD_VER=v3.3.9
curl -L ${GITHUB_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /usr/bin --strip-components=1

// Setup etcd env variable to connect flannel to it
ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_INITIAL_ADVERTISE_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCDCTL_API=3 /usr/bin/etcdctl version

/usr/bin/etcdctl set /coreos.com/network/config '{"Network":"{{ .Network }}", "Backend": {"Type": "{{ .NetworkType }}"}}'
/usr/bin/etcdctl get /coreos.com/network/config

cat << EOF > /etc/systemd/system/flanneld.service
[Unit]
Description=Networking service
Requires=etcd.service

[Service]
Restart=always

Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
Environment="ETCD_IMAGE_TAG=v3.3.9"
Environment="ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379"
Environment="ETCD_INITIAL_ADVERTISE_PEER_URLS=http://{{ .EtcdHost }}:2380"
Environment="ETCD_LISTEN_CLIENT_URLS=http://{{ .EtcdHost }}:2379"
Environment="ETCD_LISTEN_PEER_URLS=http://{{ .EtcdHost }}:2380"
Environment="ETCDCTL_API=3"
ExecStart=/usr/bin/flanneld --etcd-endpoints=http://{{ .EtcdHost }}:2379

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable flanneld.service
systemctl start flanneld.service
