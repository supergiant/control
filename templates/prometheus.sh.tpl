{{if not .RBACEnabled }}
wget https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl apply -f cluster-roles.yaml --validate=false
{{end}}
sudo /opt/bin/helm install stable/prometheus-operator