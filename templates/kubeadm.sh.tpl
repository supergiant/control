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

{{ if .IsBootstrap }}
sudo kubeadm init --token={{ .Token }} --pod-network-cidr={{ .CIDR }}
sudo kubeadm config view > kubeadm-config.yaml
sed -i 's/controlPlaneEndpoint: ""/controlPlaneEndpoint: "https://{{ .LoadBalancerHost }}"/g' kubeadm-config.yaml
sudo kubeadm config upload from-file --config=kubeadm-config.yaml

{{ if eq .NeworkProvider "Flannel" }}
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/bc79dd1505b0c8681ece4de4c0d86c5cd2643275/Documentation/kube-flannel.yml
{{ end }}


{{ if eq .NeworkProvider "Calico" }}
kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
{{ end }}

{{ if eq .NeworkProvider "Weave" }}
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
{{ end }}

{{ else }}
sudo kubeadm join https://{{ .LoadBalancerHost }} --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification --experimental-control-plane
{{ end }}

sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
{{ else }}
sudo kubeadm join https://{{ .LoadBalancerHost }} --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification
{{ end }}