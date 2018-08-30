mkdir -p /tmp/etcd-data
cat > /etc/systemd/system/etcd.service <<EOF
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
            --volume /etc/ssl/certs/:/etc/ssl/certs \
            --name {{ .Name }} \
            gcr.io/etcd-development/etcd:v{{ .Version }} \
            /usr/local/bin/etcd \
            --name {{ .Name }} \
            --data-dir /etcd-data \
            --listen-client-urls http://{{ .Host }}:{{ .ServicePort }} \
            --advertise-client-urls http://{{ .Host }}:{{ .ServicePort }} \
            --listen-peer-urls http://{{ .Host }}:{{ .ManagementPort }} \
            --initial-advertise-peer-urls http://{{ .Host }}:{{ .ManagementPort }} \
            {{if gt .ClusterSize 1 }}
            --discovery {{ .DiscoveryUrl }} \
            {{else}}
            --initial-cluster {{ .Name }}=http://{{ .Host }}:{{ .ManagementPort }} \
            {{end}}
            --listen-peer-urls http://{{ .Host }}:2380 --listen-client-urls http://{{ .Host }}:2379

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable etcd.service
systemctl start etcd.service