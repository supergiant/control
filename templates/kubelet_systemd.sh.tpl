    cat << EOF > /etc/systemd/system/{{ .KubeletService }}
[Unit]
Description=Kubernetes Kubelet Server
Documentation=https://github.com/kubernetes/kubernetes
Requires=docker.service network-online.target
After=docker.service network-online.target

[Service]
ExecStartPre=/bin/bash -c "/opt/bin/download-k8s-binary"
ExecStartPost=/bin/bash -c "/opt/bin/kube-post-start.sh"

ExecStart=/usr/bin/docker run \
        --net=host \
        --pid=host \
        --privileged \
        -v /dev:/dev \
        -v /sys:/sys:ro \
        -v /var/run:/var/run:rw \
        -v /var/lib/docker/:/var/lib/docker:rw \
        -v /var/lib/kubelet/:/var/lib/kubelet:shared \
        -v /var/log:/var/log:shared \
        -v /srv/kubernetes:/srv/kubernetes:ro \
        -v /etc/kubernetes:/etc/kubernetes:ro \
        gcr.io/google-containers/hyperkube:v{{ .K8SVersion }} \
        /hyperkube kubelet --allow-privileged=true \
        --cluster-dns=10.3.0.10 \
        --cluster_domain=cluster.local \
        --cadvisor-port=0 \
        --pod-manifest-path=/etc/kubernetes/manifests \
        --kubeconfig=/var/lib/kubelet/kubeconfig \
        --volume-plugin-dir=/etc/kubernetes/volumeplugins \
        {{- .KubernetesProvider }}
        --register-node=false
Restart=always
StartLimitInterval=0
RestartSec=10
KillMode=process

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable {{ .KubeletService }}
systemctl start {{ .KubeletService }}
