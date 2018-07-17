#!/bin/sh

# Define local vars
KUBERNETES_CONFIG_DIR=/etc/kubernetes
KUBELET_SERVICE='kubelet.service'

# Define template vars
PRIVATE_IPV4="{{ .PrivateIPV4 }}"

write_ssl_certs() {
    KUBERNETES_SSL_DIR=${KUBERNETES_CONFIG_DIR}/ssl

    mkdir -p ${KUBERNETES_SSL_DIR}
    echo "{{ .CACert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CACertName }}'
    echo "{{ .CAKeyCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CAKeyName }}'
    echo "{{ .APIServerCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerCertName }}'
    echo "{{ .APIServerKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerKeyName }}'
    echo "{{ .KubeletClientCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientCertName }}'
    echo "{{ .KubeletClientKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientKeyName }}'
}

write_manifests() {
    KUBERNETES_MANIFESTS_DIR=${KUBERNETES_CONFIG_DIR}/manifests

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
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - apiserver
    - --bind-address=0.0.0.0
    - --etcd-servers=http://127.0.0.1:2379
    - --allow-privileged=true
    {{if .RBACEnabled }}- --authorization-mode=Node,RBAC{{end}}
    - --service-cluster-ip-range=10.3.0.0/24
    - --secure-port=443
    - --v=2
    - --advertise-address=${PRIVATE_IPV4}
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
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - controller-manager
    - --master=http://127.0.0.1:8080
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
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - scheduler
    - --v=2
    - --master=http://127.0.0.1:8080
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
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - proxy
    - --v=2
    - --master=http://127.0.0.1:8080
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
}

write_addons() {
    KUBERNETES_ADDONS_DIR=${KUBERNETES_CONFIG_DIR}/addons

    ADDON=${KUBERNETES_ADDONS_DIR}/'cluster-monitoring'
    mkdir -p ${ADDON}
    cat << EOF > ${ADDON}/heapster-controller.yaml
apiVersion: v1
kind: ReplicationController
metadata:
  name: heapster-v11
  namespace: kube-system
  labels:
    k8s-app: heapster
    version: v11
    kubernetes.io/cluster-service: "true"
spec:
  replicas: 1
  selector:
    k8s-app: heapster
    version: v11
  template:
    metadata:
      labels:
        k8s-app: heapster
        version: v11
        kubernetes.io/cluster-service: "true"
    spec:
      containers:
        - image: gcr.io/google_containers/heapster:{{ .HeapsterVersion }}
          name: heapster
          resources:
            limits:
              cpu: 100m
              memory: 212Mi
          command:
            - /heapster
            - --source=kubernetes
            - --sink=influxdb:http://monitoring-influxdb:8086
            - --metric_resolution={{ .HeapsterMetricResolution }}
EOF

    cat << EOF > ${ADDON}/heapster-service.yaml
kind: Service
apiVersion: v1
metadata:
  name: heapster
  namespace: kube-system
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "Heapster"
spec:
  ports:
    - port: 80
      targetPort: 8082
  selector:
    k8s-app: heapster
EOF

    cat << EOF > ${ADDON}/influxdb-grafana-controller.yaml
apiVersion: v1
kind: ReplicationController
metadata:
  name: monitoring-influxdb-grafana-v2
  namespace: kube-system
  labels:
    k8s-app: influxGrafana
    version: v2
    kubernetes.io/cluster-service: "true"
spec:
  replicas: 1
  selector:
    k8s-app: influxGrafana
    version: v2
  template:
    metadata:
      labels:
        k8s-app: influxGrafana
        version: v2
        kubernetes.io/cluster-service: "true"
    spec:
      containers:
        - image: gcr.io/google_containers/heapster_influxdb:v0.4
          name: influxdb
          resources:
            limits:
              cpu: 100m
              memory: 200Mi
          ports:
            - containerPort: 8083
              hostPort: 8083
            - containerPort: 8086
              hostPort: 8086
          volumeMounts:
          - name: influxdb-persistent-storage
            mountPath: /data
        - image: beta.gcr.io/google_containers/heapster_grafana:v2.1.1
          name: grafana
          env:
          resources:
            limits:
              cpu: 100m
              memory: 100Mi
          env:
            # This variable is required to setup templates in Grafana.
            - name: INFLUXDB_SERVICE_URL
              value: http://monitoring-influxdb:8086
              # The following env variables are required to make Grafana accessible via
              # the kubernetes api-server proxy. On production clusters, we recommend
              # removing these env variables, setup auth for grafana, and expose the grafana
              # service using a LoadBalancer or a public IP.
            - name: GF_AUTH_BASIC_ENABLED
              value: "false"
            - name: GF_AUTH_ANONYMOUS_ENABLED
              value: "true"
            - name: GF_AUTH_ANONYMOUS_ORG_ROLE
              value: Admin
            - name: GF_SERVER_ROOT_URL
              value: /api/v1/proxy/namespaces/kube-system/services/monitoring-grafana/
          volumeMounts:
          - name: grafana-persistent-storage
            mountPath: /var
      volumes:
      - name: influxdb-persistent-storage
        emptyDir: {}
      - name: grafana-persistent-storage
        emptyDir: {}
EOF

    cat << EOF > ${ADDON}/influxdb-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: monitoring-influxdb
  namespace: kube-system
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "InfluxDB"
spec:
  ports:
    - name: http
      port: 8083
      targetPort: 8083
    - name: api
      port: 8086
      targetPort: 8086
  selector:
    k8s-app: influxGrafana
EOF

    ADDON=${KUBERNETES_ADDONS_DIR}/'kube-dns'
    mkdir -p ${ADDON}
    cat << EOF > ${ADDON}/kube-dns.yaml
apiVersion: v1
kind: Service
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "KubeDNS"
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
    kubernetes.io/cluster-service: "true"
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
EOF
}

write_kubelet_config() {
    cat << EOF > /var/lib/kubelet/kubeconfig
apiVersion: v1
kind: Config
users:
- name: kubelet
  user:
    client-certificate: /home/unknown/.minikube/client.crt
    client-key: /home/unknown/.minikube/client.key
clusters:
- name: local
  cluster:
    server: http://127.0.0.1:8080
    insecure-skip-tls-verify: true
contexts:
- name: kubelet-local
  context:
    cluster: local
    user: kubelet
current-context: kubelet-local
EOF
}

write_systemd_units() {
    cat << EOF > /etc/systemd/system/${KUBELET_SERVICE}
[Unit]
Description=Kubernetes Kubelet Server
Documentation=https://github.com/kubernetes/kubernetes
Requires=docker.service network-online.target
After=docker.service network-online.target

[Service]
ExecStartPre=/bin/bash -c "/opt/bin/download-k8s-binary"
ExecStartPost=/bin/bash -c "/opt/bin/kube-post-start.sh"

ExecStart=/usr/bin/docker run \
        --net=host \
        --pid=host \
        --privileged \
        -v /dev:/dev \
        -v /sys:/sys:ro \
        -v /var/run:/var/run:rw \
        -v /var/lib/docker/:/var/lib/docker:rw \
        -v /var/lib/kubelet/:/var/lib/kubelet:shared \
        -v /var/log:/var/log:shared \
        -v /srv/kubernetes:/srv/kubernetes:ro \
        -v /etc/kubernetes:/etc/kubernetes:ro \
        gcr.io/google-containers/hyperkube:v{{ .KubernetesVersion }} \
        /hyperkube kubelet --allow-privileged=true \
        --cluster-dns=10.3.0.10 \
        --cluster_domain=cluster.local \
        --cadvisor-port=0 \
        --pod-manifest-path=/etc/kubernetes/manifests \
        --kubeconfig=/var/lib/kubelet/kubeconfig \
        --volume-plugin-dir=/etc/kubernetes/volumeplugins \
        {{- .KubeProviderString }}
        --register-node=false
Restart=always
StartLimitInterval=0
RestartSec=10
KillMode=process

[Install]
WantedBy=multi-user.target
EOF
}

## Run

write_ssl_certs
write_manifests
write_addons
write_kubelet_config
write_systemd_units

systemctl daemon-reload
systemctl enable ${KUBELET_SERVICE}
systemctl start ${KUBELET_SERVICE}

until $(curl --output /dev/null --silent --head --fail http://127.0.0.1:8080); do printf '.'; sleep 5; done
curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://127.0.0.1:8080/api/v1/namespaces
/opt/bin/kubectl config set-cluster default-cluster --server="127.0.0.1:8080"
/opt/bin/kubectl config set-context default-system --cluster=default-cluster --user=default-admin
/opt/bin/kubectl config use-context default-system

{{if .RBACEnabled }}
/opt/bin/kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
/opt/bin/kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
/opt/bin/kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
/opt/bin/kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
/opt/bin/kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
{{end}}

/opt/bin/kubectl create -f /etc/kubernetes/addons/kube-dns.yaml
/opt/bin/kubectl create -f /etc/kubernetes/addons/cluster-monitoring
/opt/bin/kubectl create -f /etc/kubernetes/addons/default-storage-class.yaml
