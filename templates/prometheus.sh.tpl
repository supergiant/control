# prometheus-operator hac rbac enabled by default so set
# it explicitly
sudo /opt/bin/helm install stable/prometheus-operator \
    --name=prometheus-operator \
    --set global.rbac.create={{ .RBACEnabled }} \
    --set grafana.rbac.create={{ .RBACEnabled }} \
    --set kube-state-metrics.rbac.create={{ .RBACEnabled }} \
    --set prometheus-node-exporter.rbac.create={{ .RBACEnabled }}
