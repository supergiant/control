package templates

const installApp = `
set -x
sudo bash -c "cat > override.yaml <<EOF
{{ .Values }}
EOF"

sudo helm install {{ .ChartRef }} --name {{ .Name }} --namespace {{ .Namespace }} -f override.yaml --debug 
`
