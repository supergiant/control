package templates

const installApp = `
set -x
sudo bash -c "cat > override.yaml <<EOF
{{ .Values }}
EOF"

helm install {{ .RepoName }}/{{ .ChartName }}-{{ .ChartVersion }}.tgz --name {{ .Name }} --namespace {{ .Namespace }} -f override.yaml --debug 
`
