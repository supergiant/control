package templates

//
// dashboard hasn't migrated to the metrics server API and
// still depends on heapster
// https://github.com/kubernetes/dashboard/issues/2986
//
// probably, it's better to set read-only access
// --set rbac.clusterReadOnlyRole=true
const dashboardTpl = `
sudo /usr/bin/helm install stable/heapster \
   -n heapster \
   --namespace kube-system

sudo /usr/bin/helm install stable/kubernetes-dashboard \
   -n kubernetes-dashboard \
   --namespace kube-system \
   --set enableSkipLogin=true \
   --set enableInsecureLogin=true \
   --set rbac.clusterAdminRole=true
`
