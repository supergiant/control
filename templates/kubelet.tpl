sudo bash -c "cat > /etc/default/kubelet <<EOF
KUBELET_EXTRA_ARGS=--tls-cert-file=/etc/kubernetes/pki/kubelet.crt \
--tls-private-key-file=/etc/kubernetes/pki/kubelet.key \
{{ if eq .Provider "openstack" }}--cloud-provider={{ .Provider }} {{ end }} \
--rotate-certificates  --feature-gates=RotateKubeletClientCertificate=true
EOF"

sudo systemctl daemon-reload
sudo systemctl restart kubelet