echo '{{ .KubeletService }}' > /etc/systemd/system/kubelet.service
systemctl start kubelet