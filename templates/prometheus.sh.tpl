wget https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.12/plugin/pkg/auth/authorizer/rbac/bootstrappolicy/testdata/cluster-roles.yaml
sudo kubectl create -f cluster-roles.yaml --validate=false
TILLER_PORT=44135
TILLER_IP=$(kubectl describe pod tiller-deploy-57f988f854-v9x46 -n kube-system|grep IP| awk '{print $2}')
until $(curl --output /dev/null --silent --head --fail http://$TILLER_IP:$TILLER_PORT); do
    printf '.'
    sleep 5
done

sudo /opt/bin/helm install stable/prometheus-operator