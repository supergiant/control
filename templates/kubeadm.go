package templates

const kubeadmTpl = `
set -e

sudo apt-get update && sudo apt-get install -y apt-transport-https curl
sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo bash -c "cat << EOF > /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF"

sudo apt-get update
sudo apt-get install -y kubelet={{ .K8SVersion }}-00 kubeadm={{ .KubeadmVersion }}-00 kubectl={{ .K8SVersion }}-00 --allow-unauthenticated
sudo apt-mark hold kubelet kubeadm kubectl

sudo systemctl daemon-reload
sudo systemctl restart kubelet

HOSTNAME="$(hostname)"
{{ if eq .Provider "aws" }}
HOSTNAME="$(hostname -f)"
{{ end }}

sudo mkdir -p /etc/supergiant

{{if .IsMaster }}
{{ if .IsBootstrap }}

sudo bash -c "cat << EOF > /etc/supergiant/kubeadm.conf
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
localAPIEndpoint:
  bindPort: 443
nodeRegistration:
  kubeletExtraArgs:
    node-ip: {{ .NodeIp }}
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v{{ .K8SVersion }}
clusterName: kubernetes
controlPlaneEndpoint: {{ .InternalDNSName }}:443
certificatesDir: /etc/kubernetes/pki
apiServer:
  certSANs:
  - {{ .ExternalDNSName }}
  - {{ .InternalDNSName }}
  extraArgs:
    authorization-mode: Node,RBAC
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
  timeoutForControlPlane: 8m0s
controllerManager:
  extraArgs:
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
networking:
  dnsDomain: cluster.local
  podSubnet: {{ .CIDR }}
  serviceSubnet: {{ .ServiceCIDR }}
EOF"

sudo kubeadm init --ignore-preflight-errors=NumCPU \
--node-name ${HOSTNAME} \
--config=/etc/supergiant/kubeadm.conf \
--upload-certs
{{ else }}

sudo bash -c "cat << EOF > /etc/supergiant/kubeadm.conf
apiVersion: kubeadm.k8s.io/v1beta1
kind: JoinConfiguration
nodeRegistration:
  kubeletExtraArgs:
    node-ip: {{ .NodeIp }}
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
discovery:
  bootstrapToken:
    token: {{ .Token }}
    apiServerEndpoint: {{ .InternalDNSName }}:443
    unsafeSkipCAVerification: true
controlPlane:
  localAPIEndpoint:
    bindPort: 443
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v{{ .K8SVersion }}
clusterName: kubernetes
controlPlaneEndpoint: {{ .InternalDNSName }}:443
certificatesDir: /etc/kubernetes/pki
apiServer:
  certSANs:
  - {{ .ExternalDNSName }}
  - {{ .InternalDNSName }}
  extraArgs:
    authorization-mode: Node,RBAC
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
  timeoutForControlPlane: 8m0s
controllerManager:
  extraArgs:
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
networking:
  dnsDomain: cluster.local
  podSubnet: {{ .CIDR }}
  serviceSubnet: {{ .ServiceCIDR }}
EOF"

sudo kubeadm config images pull
sudo kubeadm join --ignore-preflight-errors=NumCPU {{ .InternalDNSName }}:443 \
--node-name ${HOSTNAME} \
--certificate-key {{ .CertificateKey }} \
--config=/etc/supergiant/kubeadm.conf
{{ end }}

sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
{{ else }}

sudo bash -c "cat << EOF > /etc/supergiant/kubeadm.conf
apiVersion: kubeadm.k8s.io/v1beta1
kind: JoinConfiguration
nodeRegistration:
  kubeletExtraArgs:
    node-ip: {{ .NodeIp }}
    {{ if .Provider }}cloud-provider: {{ .Provider }}{{ end }}
discovery:
  bootstrapToken:
    token: {{ .Token }}
    apiServerEndpoint: {{ .InternalDNSName }}:443
    unsafeSkipCAVerification: true
EOF"

sudo kubeadm join --ignore-preflight-errors=NumCPU {{ .InternalDNSName }}:443 \
--node-name ${HOSTNAME} \
--config=/etc/supergiant/kubeadm.conf
{{ end }}
`
