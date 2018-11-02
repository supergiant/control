{{if not .RBACEnabled }}
wget https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl apply -f cluster-roles.yaml --validate=false
{{end}}
sleep 60
sudo /opt/bin/helm install stable/prometheus-operator

sudo bash -c "cat << EOF > prometheus.yaml
apiVersion: v1
kind: Service
metadata:
  name: prometheus
spec:
  type: NodePort
  ports:
  - name: web
    nodePort: {{ .Port }}
    port: 9090
    protocol: TCP
    targetPort: web
  selector:
    app: prometheus
EOF"

sudo kubectl create -f prometheus.yaml