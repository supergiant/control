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