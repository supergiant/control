package templates

const upgradeTpl = `
sudo apt update
sudo apt-cache policy kubeadm

sudo apt-mark unhold kubeadm && \
sudo apt-get update && sudo apt-get install -y kubeadm={{ .K8SVersion }}-00 && \
sudo apt-mark hold kubeadm


{{ if .IsBootstrap }}
	sudo kubeadm upgrade plan
	sudo kubeadm upgrade apply -y v{{ .K8SVersion }}
{{ else }}
	{{ if .IsMaster }}
		sudo kubeadm upgrade node experimental-control-plane
	{{ else }}
		sudo kubeadm upgrade node config --kubelet-version v{{ .K8SVersion }}
	{{ end }}
{{ end }}

sudo apt-mark unhold kubelet kubectl && \
sudo apt-get update && sudo apt-get install -y kubelet={{ .K8SVersion }}-00 kubectl={{ .K8SVersion }}-00 && \
sudo apt-mark hold kubelet kubectl
sudo systemctl restart kubelet
`
