sudo mkdir -p {{ .DataDir }}

ETCD_VER=v3.3.9

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/coreos/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

sudo rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
sudo rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

sudo curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
sudo tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /usr/bin --strip-components=1
sudo rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

sudo etcd --version
sudo ETCDCTL_API=3 etcdctl version

sudo bash -c "cat > /etc/systemd/system/etcd.service <<EOF
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd
Requires=network-online.target

[Service]
Restart=always
RestartSec={{ .RestartTimeout }}s
LimitNOFILE=40000
TimeoutStartSec={{ .StartTimeout }}s

ExecStart=/usr/bin/etcd --name {{ .Name }} \
            --data-dir {{ .DataDir }} \
            --listen-client-urls http://{{ .Host }}:{{ .ServicePort }} \
            --advertise-client-urls http://{{ .AdvertiseHost }}:{{ .ServicePort }} \
            --listen-peer-urls http://{{ .Host }}:{{ .ManagementPort }} \
            --initial-advertise-peer-urls http://{{ .AdvertiseHost }}:{{ .ManagementPort }} \
            --discovery {{ .DiscoveryUrl }} \

[Install]
WantedBy=multi-user.target
EOF"
sudo systemctl daemon-reload
sudo systemctl enable etcd.service
sudo systemctl start etcd.service

while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://{{ .Host }}:{{ .ServicePort }}/health)" != "200" ]]; do printf 'wait for etcd\n';sleep 5; done
