sudo apt-get update
sudo wget https://dl.k8s.io/v{{ .K8SVersion }}/kubernetes-server-linux-amd64.tar.gz
sudo tar -xvf kubernetes-server-linux-amd64.tar.gz

sudo cp kubernetes/server/bin/kubelet /usr/bin
sudo cp kubernetes/server/bin/kubectl /usr/bin
sudo cp kubernetes/server/bin/kubeadm /usr/bin

sudo bash -c "cat > /etc/systemd/system/kubelet.service <<EOF
[Service]
Environment=KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf
Environment=KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
EnvironmentFile=-/etc/default/kubelet
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
EOF"

sudo bash -c "cat > /etc/default/kubelet <<EOF
KUBELET_EXTRA_ARGS=--tls-cert-file=/etc/kubernetes/pki/kubelet.crt \
--tls-private-key-file=/etc/kubernetes/pki/kubelet.key \
{{ if eq .Provider "openstack" }}--cloud-provider={{ .Provider }} {{ end }} \
--rotate-certificates  --feature-gates=RotateKubeletClientCertificate=true
EOF"

sudo systemctl daemon-reload
sudo systemctl restart kubelet