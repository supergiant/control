wget https://github.com/kubernetes/kubernetes/blob/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl create -f cluster-roles.yaml --validate=false
/opt/bin/helm install stable/prometheus-operator