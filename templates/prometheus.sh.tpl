# prometheus-operator hac rbac enabled by default so set
# it explicitly
sudo /opt/bin/helm install stable/prometheus-operator \
    --name=prometheus-operator \
    --set prometheus.prometheusSpec.externalUrl=/api/v1/namespaces/default/services/prometheus-operated:9090/proxy/ \
        --set alertmanager.alertmanagerSpec.externalUrl=/api/v1/namespaces/default/services/prometheus-operated:9090/proxy/ \
    --set global.rbac.create={{ .RBACEnabled }} \
    --set grafana.rbac.create={{ .RBACEnabled }} \
    --set kube-state-metrics.rbac.create={{ .RBACEnabled }} \
    --set prometheus-node-exporter.rbac.create={{ .RBACEnabled }}