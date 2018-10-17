KUBERNETES_SSL_DIR={{ .KubernetesConfigDir }}/ssl

mkdir -p ${KUBERNETES_SSL_DIR}

FLANNELIP=$(ifconfig flannel.1|grep 'inet addr'| awk '{print substr($2,6,12)}')

cat > /etc/kubernetes/ssl/openssl.cnf.template <<EOF
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
IP.1 = {FLANNELIP}
IP.2 = {MASTER_HOST}
IP.3 = {PRIVATE_HOST}
EOF

cat > /etc/kubernetes/ssl/ca.pem <<EOF
{{ .CACert }}EOF

cat > /etc/kubernetes/ssl/ca-key.pem <<EOF
{{ .CAKey }}EOF

sed -e "s/{MASTER_HOST}/`curl ipinfo.io/ip`/" < /etc/kubernetes/ssl/openssl.cnf.template > /etc/kubernetes/ssl/openssl.cnf.1
sed -e "s/{PRIVATE_HOST}/{{ .MasterHost }}/" < /etc/kubernetes/ssl/openssl.cnf.1 > /etc/kubernetes/ssl/openssl.cnf.2
sed -e "s/{FLANNELIP}/$FLANNELIP/" < /etc/kubernetes/ssl/openssl.cnf.2 > /etc/kubernetes/ssl/openssl.cnf

openssl genrsa -out /etc/kubernetes/ssl/apiserver-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/apiserver-key.pem -out /etc/kubernetes/ssl/apiserver.csr -subj "/CN=kube-apiserver" -config /etc/kubernetes/ssl/openssl.cnf
openssl x509 -req -in /etc/kubernetes/ssl/apiserver.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/apiserver.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
openssl genrsa -out /etc/kubernetes/ssl/worker-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/worker-key.pem -out /etc/kubernetes/ssl/worker.csr -subj "/CN=kube-worker"
openssl x509 -req -in /etc/kubernetes/ssl/worker.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/worker.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
openssl genrsa -out /etc/kubernetes/ssl/admin-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/admin-key.pem -out /etc/kubernetes/ssl/admin.csr -subj "/CN=kube-admin"
openssl x509 -req -in /etc/kubernetes/ssl/admin.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/admin.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf

cp /etc/kubernetes/ssl/ca.pem /usr/share/ca-certificates/ca.crt
echo "ca.crt" >> /etc/ca-certificates.conf
update-ca-certificates

chmod 600 /etc/kubernetes/ssl/*-key.pem
chown root:root /etc/kubernetes/ssl/*-key.pem

cat > /etc/kubernetes/ssl/basic_auth.csv <<EOF
{{ .Password }},{{ .Username }},admin
EOF

cat > /etc/kubernetes/ssl/known_tokens.csv <<EOF
{{ .Password }},kubelet,kubelet
{{ .Password }},kube_proxy,kube_proxy
{{ .Password }},system:scheduler,system:scheduler
{{ .Password }},system:controller_manager,system:controller_manager
{{ .Password }},system:logging,system:logging
{{ .Password }},system:monitoring,system:monitoring
{{ .Password }},system:dns,system:dns
EOF