package templates

const cloudcontrollerTpl = `
sudo bash -c 'cat << EOF | kubectl create -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin ## TODO: add cloud-controller-manager role
subjects:
- kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
{{- if eq .Provider "digitalocean" }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: digitalocean-cloud-controller-manager
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 2
  selector:
    matchLabels:
      app: digitalocean-cloud-controller-manager
  template:
    metadata:
      labels:
        app: digitalocean-cloud-controller-manager
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
    spec:
      dnsPolicy: Default
      hostNetwork: true
      serviceAccountName: cloud-controller-manager
      tolerations:
        - key: "node.cloudprovider.kubernetes.io/uninitialized"
          value: "true"
          effect: "NoSchedule"
        - key: "CriticalAddonsOnly"
          operator: "Exists"
        - key: "node-role.kubernetes.io/master"
          effect: NoSchedule
      containers:
      - image: digitalocean/digitalocean-cloud-controller-manager:v0.1.9
        name: digitalocean-cloud-controller-manager
        command:
          - "/bin/digitalocean-cloud-controller-manager"
          - "--cloud-provider=digitalocean"
          - "--leader-elect=false"
        resources:
          requests:
            cpu: 100m
            memory: 50Mi
        env:
          - name: DO_ACCESS_TOKEN # TODO: use secrets
            value: "{{ .DOAccessToken }}"
{{ end }}
EOF'
`
