wget https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl create -f cluster-roles.yaml --validate=false
until $([ $(sudo kubectl get pods --namespace=kube-system|grep tiller|grep Ready|wc -l) -eq 1 ]); do printf '.'; sleep 5; done
sudo /opt/bin/helm install stable/prometheus-operator