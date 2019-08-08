package templates

const installApp = `
set -x
sudo bash -c "cat > override.yaml <<EOF
{{ .Values }}
EOF"

sudo helm install {{ .ChartRef }} {{ if .Name }}--name {{ .Name }}{{ end}} --namespace {{ .Namespace }} -f override.yaml --debug 
`
