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
            --name {{ .Name }} \
            gcr.io/etcd-development/etcd:v{{ .Version }} \
            /usr/local/bin/etcd \
            --name {{ .Name }} \
            --data-dir /etcd-data \
            --listen-client-urls http://{{ .MasterPrivateIP }}:{{ .ServicePort }} \
            --advertise-client-urls http://{{ .MasterPrivateIP }}:{{ .ServicePort }} \
            --listen-peer-urls http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --initial-advertise-peer-urls http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --initial-cluster {{ .Name }}=http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --listen-peer-urls http://{{ .MasterPrivateIP }}:2380 --listen-client-urls http://{{ .MasterPrivateIP }}:2379

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable etcd.service
systemctl start etcd.service