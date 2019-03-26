{{if .IsMaster }}
sudo kubeadm config images pull

{{ if .IsBootstrap }}
sudo kubeadm init --ignore-preflight-errors=NumCPU,Port-10250 --apiserver-advertise-address={{ .AdvertiseAddress }} --token={{ .Token }} --pod-network-cidr={{ .CIDR }} \
--kubernetes-version {{ .K8SVersion }} --apiserver-bind-port=443 --apiserver-cert-extra-sans {{ .ExternalDNSName }},{{ .InternalDNSName }}
sudo kubeadm config view > kubeadm-config.yaml
sed -i 's/controlPlaneEndpoint: ""/controlPlaneEndpoint: "{{ .InternalDNSName }}:443"/g' kubeadm-config.yaml
sudo kubeadm config upload from-file --config=kubeadm-config.yaml

{{ else }}
sudo kubeadm join --ignore-preflight-errors=NumCPU,Port-10250 {{ .InternalDNSName }}:443 --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification --experimental-control-plane
{{ end }}

sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
{{ else }}
sudo kubeadm join --ignore-preflight-errors=NumCPU,Port-10250 {{ .InternalDNSName }}:443 --token {{ .Token }} \
--discovery-token-unsafe-skip-ca-verification
{{ end }}