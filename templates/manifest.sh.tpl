KUBERNETES_MANIFESTS_DIR={{ .KubernetesConfigDir }}/manifests
KUBERNETES_ADDONS_DIR={{ .KubernetesConfigDir }}/addons

ADDON=${KUBERNETES_ADDONS_DIR}/'kube-dns'
sudo mkdir -p ${ADDON}
sudo bash -c "cat << EOF > ${ADDON}/kube-dns.yaml
apiVersion: v1
kind: Service
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    kubernetes.io/cluster-service: 'true'
    kubernetes.io/name: 'KubeDNS'
spec:
  selector:
    k8s-app: kube-dns
  clusterIP: 10.3.0.10
  ports:
  - name: dns
    port: 53
    protocol: UDP
  - name: dns-tcp
    port: 53
    protocol: TCP
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: kube-dns-v11
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    version: v11
    kubernetes.io/cluster-service: 'true'
spec:
  replicas: 1
  selector:
    k8s-app: kube-dns
    version: v11
  template:
    metadata:
      labels:
        k8s-app: kube-dns
        version: v11
        kubernetes.io/cluster-service: 'true'
    spec:
        replicas: 1
        selector:
          k8s-app: kube-dns
          version: v11
        template:
          metadata:
            labels:
              k8s-app: kube-dns
              version: v11
              kubernetes.io/cluster-service: "true"
          spec:
            containers:
            - name: etcd
              image: gcr.io/google_containers/etcd:2.0.9
              resources:
                limits:
                  cpu: 100m
                  memory: 50Mi
              command:
              - /usr/local/bin/etcd
              - -data-dir
              - /var/etcd/data
              - -listen-client-urls
              - http://127.0.0.1:2379,http://127.0.0.1:4001
              - -advertise-client-urls
              - http://127.0.0.1:2379,http://127.0.0.1:4001
              - -initial-cluster-token
              - skydns-etcd
              volumeMounts:
              - name: etcd-storage
                mountPath: /var/etcd/data
            - name: kube2sky
              image: gcr.io/google_containers/kube2sky:1.11
              resources:
                limits:
                  cpu: 100m
                  memory: 50Mi
              args:
              # command = "/kube2sky"
              - -domain=cluster.local
            - name: skydns
              image: gcr.io/google_containers/skydns:2015-03-11-001
              resources:
                limits:
                  cpu: 100m
                  memory: 50Mi
              args:
              # command = "/skydns"
              - -machines=http://localhost:4001
              - -addr=0.0.0.0:53
              - -domain=cluster.local.
              ports:
              - containerPort: 53
                name: dns
                protocol: UDP
              - containerPort: 53
                name: dns-tcp
                protocol: TCP
              livenessProbe:
                httpGet:
                  path: /healthz
                  port: 8080
                  scheme: HTTP
                initialDelaySeconds: 30
                timeoutSeconds: 5
              readinessProbe:
                httpGet:
                  path: /healthz
                  port: 8080
                  scheme: HTTP
                initialDelaySeconds: 1
                timeoutSeconds: 5
            - name: healthz
              image: gcr.io/google_containers/exechealthz:1.0
              resources:
                limits:
                  cpu: 10m
                  memory: 20Mi
              args:
              - -cmd=nslookup kubernetes.default.svc.cluster.local localhost >/dev/null
              - -port=8080
              ports:
              - containerPort: 8080
                protocol: TCP
            volumes:
            - name: etcd-storage
              emptyDir: {}
            dnsPolicy: Default
EOF"

sudo mkdir -p ${KUBERNETES_MANIFESTS_DIR}

# worker
sudo bash -c "cat << EOF > {{ .KubernetesConfigDir }}/worker-kubeconfig.yaml
apiVersion: v1
kind: Config
users:
- name: kubelet
  user:
    token: '1234'
clusters:
- name: local
  cluster:
    insecure-skip-tls-verify: true
    server: https://{{ .MasterHost }}
contexts:
- context:
    cluster: local
    user: kubelet
  name: service-account-context
current-context: service-account-context
EOF"


# proxy
sudo bash -c "cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-proxy.yaml
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
EOF"


{{ if .IsMaster }}
# api-server
sudo bash -c "cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-apiserver.yaml
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
    - --etcd-servers=http://{{ .MasterHost }}:2379
    - --allow-privileged=true
    - --service-cluster-ip-range=10.3.0.0/24
    - --secure-port=443
    - --v=2
    - --insecure-port=8080
    - --insecure-bind-address=0.0.0.0
    {{if .RBACEnabled }}- --authorization-mode=Node,RBAC{{end}}
    - --advertise-address={{ .MasterHost }}
    - --admission-control=NamespaceLifecycle,NamespaceExists,LimitRanger,ServiceAccount,ResourceQuota,DefaultStorageClass{{if .RBACEnabled }},NodeRestriction{{end}}
    - --tls-cert-file=/etc/kubernetes/ssl/apiserver.pem
    - --tls-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --kubelet-certificate-authority=/etc/kubernetes/ssl/ca.pem
    - --kubelet-client-certificate=/etc/kubernetes/ssl/apiserver.pem
    - --kubelet-client-key=/etc/kubernetes/ssl/apiserver-key.pem
    - --client-ca-file=/etc/kubernetes/ssl/ca.pem
    - --service-account-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --basic-auth-file=/etc/kubernetes/ssl/basic_auth.csv
    - --token-auth-file=/etc/kubernetes/ssl/known_tokens.csv
    - --kubelet-preferred-address-types=InternalIP,Hostname,ExternalIP
    - --storage-backend=etcd3
    -  {{ .ProviderString }}
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
EOF"

# kube controller manager
sudo bash -c "cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-controller-manager.yaml
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
    -  {{ .ProviderString }}
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
EOF"

# scheduler
sudo bash -c "cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-scheduler.yaml
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
EOF"
{{ end }}