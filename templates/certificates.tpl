sudo mkdir -p /etc/kubernetes
sudo mkdir -p /etc/kubernetes/pki
sudo mkdir -p /etc/kubernetes/pki/etcd

sudo bash -c "cat > /etc/kubernetes/pki/openssl.cnf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster
DNS.5 = kubernetes.default.svc.cluster.local
IP.1 = {{ .PublicIP }}
IP.2 = {{ .PrivateIP }}
IP.3 = {{ .KubernetesSvcIP }}
EOF"


{{ if .IsMaster }}
sudo bash -c "cat > /etc/kubernetes/pki/ca.crt <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/ca.key <<EOF
{{ .CAKey }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/front-proxy-ca.crt <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/front-proxy-ca.key <<EOF
{{ .CAKey }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/etcd/ca.crt <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/etcd/ca.key <<EOF
{{ .CAKey }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/cluster-ca.crt <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/cluster-ca.key <<EOF
{{ .CAKey }}EOF"

sudo openssl genrsa -out /etc/kubernetes/pki/kubelet.key 2048
sudo openssl req -new -key /etc/kubernetes/pki/kubelet.key -out /etc/kubernetes/pki/kubelet.csr -subj "/CN=kube-worker"
sudo openssl x509 -req -in /etc/kubernetes/pki/kubelet.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial -out /etc/kubernetes/pki/kubelet.crt -days 365 -extensions v3_req -extfile /etc/kubernetes/pki/openssl.cnf

{{ else }}

sudo bash -c "cat > /etc/kubernetes/pki/ca.crt <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/admin.crt <<EOF
{{ .AdminCert }}EOF"

sudo bash -c "cat > /etc/kubernetes/pki/admin.key <<EOF
{{ .AdminKey }}EOF"

sudo kubectl --kubeconfig=/root/.kube/config config set-cluster kubernetes --server='https://{{ .MasterHost }}:{{ .MasterPort }}' --certificate-authority=/etc/kubernetes/pki/ca.crt --embed-certs=true
sudo kubectl --kubeconfig=/root/.kube/config config set-credentials kubernetes --client-certificate=/etc/kubernetes/pki/admin.crt --client-key=/etc/kubernetes/pki/admin.key --embed-certs=true
sudo kubectl --kubeconfig=/root/.kube/config config set-context kubernetes --cluster=kubernetes --user=kubernetes
sudo kubectl --kubeconfig=/root/.kube/config config use-context kubernetes

sudo openssl genrsa -out /etc/kubernetes/pki/kubelet.key 2048
sudo openssl req -new -key /etc/kubernetes/pki/kubelet.key -out /etc/kubernetes/pki/kubelet.csr -subj "/CN=kube-worker"

sudo bash -c "cat > /etc/kubernetes/pki/request.yaml <<EOF
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: {{ .NodeName }}
spec:
  groups:
  - system:authenticated
  request: $(cat /etc/kubernetes/pki/kubelet.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF"

sudo kubectl --kubeconfig=/root/.kube/config create -f /etc/kubernetes/pki/request.yaml
sudo kubectl --kubeconfig=/root/.kube/config certificate approve -f /etc/kubernetes/pki/request.yaml
sudo bash -c "cat > /etc/kubernetes/pki/kubelet.crt <<EOF
$(sudo kubectl --kubeconfig=/root/.kube/config get csr {{ .NodeName }} -o jsonpath='{.status.certificate}' | base64 -d)
EOF"

sudo rm /etc/kubernetes/pki/ca.crt
sudo rm /etc/kubernetes/pki/admin.key
sudo rm /etc/kubernetes/pki/admin.crt

{{ end }}

