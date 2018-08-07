KUBERNETES_MANIFESTS_DIR={{ .KubernetesConfigDir }}/manifests

mkdir -p ${KUBERNETES_MANIFESTS_DIR}
cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-apiserver.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-apiserver
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-apiserver
    image: gcr.io/google_containers/hyperkube:v{{ .K8SVersion }}
    command:
    - /hyperkube
    - apiserver
    - --bind-address=0.0.0.0
    - --etcd-servers=http://{{ .EtcdHost }}:{{ .EtcdPort }}
    - --allow-privileged=true
    {{if .RBACEnabled }}- --authorization-mode=Node,RBAC{{end}}
    - --service-cluster-ip-range=10.3.0.0/24
    - --secure-port=443
    - --v=2
    - --advertise-address={{ .PrivateIpv4 }}
    - --admission-control=NamespaceLifecycle,NamespaceExists,LimitRanger,ServiceAccount,ResourceQuota,DefaultStorageClass{{if .RBACEnabled }},NodeRestriction{{end}}
    - --tls-cert-file=/etc/kubernetes/ssl/apiserver.pem
    - --tls-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --client-ca-file=/etc/kubernetes/ssl/ca.pem
    - --service-account-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --basic-auth-file=/etc/kubernetes/ssl/basic_auth.csv
    - --token-auth-file=/etc/kubernetes/ssl/known_tokens.csv
    - --kubelet-preferred-address-types=InternalIP,Hostname,ExternalIP
    - --storage-backend=etcd2
    {{- .ProviderString }}
    ports:
    - containerPort: 443
      hostPort: 443
      name: https
    - containerPort: 8080
      hostPort: 8080
      name: local
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/kubernetes/addons
      name: api-addons-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /etc/kubernetes/addons
    name: api-addons-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-controller-manager.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-controller-manager
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-controller-manager
    image: gcr.io/google_containers/hyperkube:v{{ .K8SVersion }}
    command:
    - /hyperkube
    - controller-manager
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    - --service-account-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --root-ca-file=/etc/kubernetes/ssl/ca.pem
    - --v=2
    - --cluster-cidr=10.244.0.0/14
    - --allocate-node-cidrs=true
    {{- .ProviderString }}
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10252
      initialDelaySeconds: 15
      timeoutSeconds: 1
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-scheduler.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-scheduler
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-scheduler
    image: gcr.io/google_containers/hyperkube:v{{ .K8SVersion }}
    command:
    - /hyperkube
    - scheduler
    - --v=2
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10251
      initialDelaySeconds: 15
      timeoutSeconds: 1
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-proxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-proxy
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/hyperkube:v{{ .K8SVersion }}
    command:
    - /hyperkube
    - proxy
    - --v=2
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    - --proxy-mode=iptables
    securityContext:
      privileged: true
    volumeMounts:
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
