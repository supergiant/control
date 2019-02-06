sudo apt-get update && apt-get install -y apt-transport-https curl
sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

sudo cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
sudo apt install -y docker.io

sudo systemctl daemon-reload
sudo systemctl enable docker
sudo systemctl restart kubelet

{{if .IsMaster }}
kubeadm config images pull

{{ if .Bootstrap }}
sudo kubeadm init --token={{ .Token }} --pod-network-cidr={{ .CIDR }}
sudo kubeadm config view > kubeadm-config.yaml
sed -i 's/controlPlaneEndpoint: ""/controlPlaneEndpoint: "{{ .LoadBalancerIP }}:{{ .LoadBalancerPort }}"/g' kubeadm-config.yaml
sudo kubeadm config upload from-file --config=kubeadm-config.yaml

{{ else }}
sudo kubeadm join {{ .LoadBalancerIP }}:{{ .LoadBalancerPort }} --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification --experimental-control-plane
{{ end }}

sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
{{ else }}
sudo kubeadm join {{ .LoadBalancerIP }}:{{ .LoadBalancerPort }} --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification
{{ end }}