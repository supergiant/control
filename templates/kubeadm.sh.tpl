sudo apt-get update && sudo apt-get install -y apt-transport-https curl
sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo bash -c "cat << EOF > /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF"

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl --allow-unauthenticated
sudo apt-mark hold kubelet kubeadm kubectl

sudo systemctl daemon-reload
sudo systemctl restart kubelet

{{if .IsMaster }}
sudo kubeadm config images pull

{{ if .IsBootstrap }}
sudo kubeadm init --ignore-preflight-errors=NumCPU --apiserver-advertise-address={{ .AdvertiseAddress }} --token={{ .Token }} --pod-network-cidr={{ .CIDR }} \
--kubernetes-version {{ .K8SVersion }} --apiserver-bind-port=443 --apiserver-cert-extra-sans {{ .InternalDNSName }},{{ .ExternalDNSName }}
sudo kubeadm config view > kubeadm-config.yaml
sed -i 's/controlPlaneEndpoint: ""/controlPlaneEndpoint: "{{ .InternalDNSName }}:443"/g' kubeadm-config.yaml
sudo kubeadm config upload from-file --config=kubeadm-config.yaml

{{ else }}
sudo kubeadm join --ignore-preflight-errors=NumCPU {{ .InternalDNSName }}:443 --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification --experimental-control-plane
{{ end }}

sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
{{ else }}
sudo kubeadm join {{ .InternalDNSName }}:443 --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification
{{ end }}