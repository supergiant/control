package templates

const installApp = `
set -x
sudo bash -c "cat > override.yaml <<EOF
{{ .Values }}
EOF"

helm install {{ .RepoName }}/{{ .ChartName }}-{{ .ChartVersion }} --name {{ .Name }} --namespace {{ .Namespace }} -f override.yaml --debug 
`
