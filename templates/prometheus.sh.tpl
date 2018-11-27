# prometheus-operator hac rbac enabled by default so set
# it explicitly
sudo /opt/bin/helm install stable/prometheus-operator \
    --name=prometheus-operator \
    --set prometheus.prometheusSpec.externalUrl=/api/v1/namespaces/default/services/prometheus-operated:9090/proxy/ \
        --set alertmanager.alertmanagerSpec.externalUrl=/api/v1/namespaces/default/services/prometheus-operated:9090/proxy/ \
    --set global.rbac.create=true \
    --set grafana.rbac.create=true \
    --set kube-state-metrics.rbac.create=true \
    --set prometheus-node-exporter.rbac.create=true