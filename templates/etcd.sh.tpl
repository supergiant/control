sudo mkdir -p {{ .DataDir }}
sudo bash -c "cat > /etc/systemd/system/etcd.service <<EOF
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd

[Service]
Restart=always
RestartSec={{ .RestartTimeout }}s
LimitNOFILE=40000
TimeoutStartSec={{ .StartTimeout }}s

ExecStart=/usr/bin/docker run \
            -p {{ .ServicePort }}:{{ .ServicePort }} \
            -p {{ .ManagementPort }}:{{ .ManagementPort }} \
            --volume={{ .DataDir }}:/etcd-data \
            --volume=/etc/ssl/certs:/etc/ssl/certs \
            gcr.io/etcd-development/etcd:v{{ .Version }} \
            /usr/local/bin/etcd \
            --name {{ .Name }} \
            --data-dir /etcd-data \
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
