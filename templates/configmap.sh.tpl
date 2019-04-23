sudo bash -c "cat > capacity_configmap.yaml <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: capacity-configmap
  namespace: {{ .Namespace }}
data:
  {{ .Data }}
EOF"
