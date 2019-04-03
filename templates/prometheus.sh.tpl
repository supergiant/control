# prometheus-operator hac rbac enabled by default so set
# it explicitly
sudo /usr/bin/helm install stable/prometheus-operator \
    --name=prometheus-operator \
    --namespace=kube-system \
    --version 5.0.4 \
    --set global.rbac.create={{ .RBACEnabled }} \
    --set grafana.rbac.create={{ .RBACEnabled }} \
    --set kube-state-metrics.rbac.create={{ .RBACEnabled }} \
    --set prometheus-node-exporter.rbac.create={{ .RBACEnabled }} \
    --set exporter-kubelets.https=true
