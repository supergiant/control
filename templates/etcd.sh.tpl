sudo mkdir -p {{ .ETCDConfig.DataDir }}

ETCD_VER=v{{ .ETCDConfig.Version}}

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
RestartSec={{ .ETCDConfig.RestartTimeout }}s
LimitNOFILE=40000
TimeoutStartSec={{ .ETCDConfig.StartTimeout }}s

ExecStart=/usr/bin/etcd --name {{ .ETCDConfig.Name }} \
            --data-dir {{ .ETCDConfig.DataDir }} \
            --listen-client-urls http://0.0.0.0:{{ .ETCDConfig.ServicePort }} \
            --advertise-client-urls http://{{ .NodePrivateIP }}:{{ .ETCDConfig.ServicePort }} \
            --listen-peer-urls http://0.0.0.0:{{ .ETCDConfig.ManagementPort }} \
            --initial-cluster-token {{ .ETCDConfig.ClusterToken }} \
            --initial-cluster {{ .InitialClusterIPs }} \
            --initial-cluster-state new \
            --initial-advertise-peer-urls http://{{ .NodePrivateIP }}:{{ .ETCDConfig.ManagementPort }}
[Install]
WantedBy=multi-user.target
EOF"
sudo systemctl daemon-reload
sudo systemctl enable etcd.service
sudo systemctl start etcd.service

while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://{{ .NodePrivateIP }}:{{ .ETCDConfig.ServicePort }}/health)" != "200" ]]; do printf 'wait for etcd\n';sleep 5; done
