wget https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl create -f cluster-roles.yaml --validate=false
sudo /opt/bin/helm install stable/prometheus-operator