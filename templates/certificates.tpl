KUBERNETES_SSL_DIR={{ .KubernetesConfigDir }}/ssl

sudo mkdir -p ${KUBERNETES_SSL_DIR}

#sleep 10
#FLANNELIP=$(ifconfig flannel.1|grep 'inet addr'| awk '{print substr($2,6,12)}')

sudo bash -c "cat > /etc/kubernetes/ssl/openssl.cnf.template <<EOF
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
IP.1 = {MASTER_HOST}
IP.2 = {PRIVATE_HOST}
IP.3 = 10.3.0.1
EOF"

sudo bash -c "cat > /etc/kubernetes/ssl/ca.pem <<EOF
{{ .CACert }}EOF"

sudo bash -c "cat > /etc/kubernetes/ssl/ca-key.pem <<EOF
{{ .CAKey }}EOF"

sudo bash -c "sed -e \"s/{MASTER_HOST}/`curl ipinfo.io/ip`/\" < /etc/kubernetes/ssl/openssl.cnf.template > /etc/kubernetes/ssl/openssl.cnf.1"
sudo bash -c "sed -e \"s/{PRIVATE_HOST}/{{ .MasterHost }}/\" < /etc/kubernetes/ssl/openssl.cnf.1 > /etc/kubernetes/ssl/openssl.cnf"

#sudo sed -e "s/{FLANNELIP}/$FLANNELIP/" < /etc/kubernetes/ssl/openssl.cnf.2 > /etc/kubernetes/ssl/openssl.cnf

sudo openssl genrsa -out /etc/kubernetes/ssl/apiserver-key.pem 2048
sudo openssl req -new -key /etc/kubernetes/ssl/apiserver-key.pem -out /etc/kubernetes/ssl/apiserver.csr -subj "/CN=kube-apiserver" -config /etc/kubernetes/ssl/openssl.cnf
sudo openssl x509 -req -in /etc/kubernetes/ssl/apiserver.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/apiserver.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
sudo openssl genrsa -out /etc/kubernetes/ssl/worker-key.pem 2048
sudo openssl req -new -key /etc/kubernetes/ssl/worker-key.pem -out /etc/kubernetes/ssl/worker.csr -subj "/CN=kube-worker"
sudo openssl x509 -req -in /etc/kubernetes/ssl/worker.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/worker.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
sudo openssl genrsa -out /etc/kubernetes/ssl/admin-key.pem 2048
sudo openssl req -new -key /etc/kubernetes/ssl/admin-key.pem -out /etc/kubernetes/ssl/admin.csr -subj "/CN=kube-admin"
sudo openssl x509 -req -in /etc/kubernetes/ssl/admin.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/admin.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf

sudo cp /etc/kubernetes/ssl/ca.pem /usr/share/ca-certificates/ca.crt
sudo bash -c "echo \"ca.crt\" >> /etc/ca-certificates.conf"
sudo update-ca-certificates

sudo chmod 600 /etc/kubernetes/ssl/*-key.pem
sudo chown root:root /etc/kubernetes/ssl/*-key.pem

sudo bash -c "cat > /etc/kubernetes/ssl/basic_auth.csv <<EOF
{{ .Password }},{{ .Username }},admin
EOF"

sudo bash -c "cat > /etc/kubernetes/ssl/known_tokens.csv <<EOF
{{ .Password }},kubelet-bootstrap,10001,"system:node-bootstrapper"
{{ .Password }},kubelet,kubelet
{{ .Password }},kube_proxy,kube_proxy
{{ .Password }},system:scheduler,system:scheduler
{{ .Password }},system:controller_manager,system:controller_manager
{{ .Password }},system:logging,system:logging
{{ .Password }},system:monitoring,system:monitoring
{{ .Password }},system:dns,system:dns
EOF"