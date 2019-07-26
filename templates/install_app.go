package templates

const installApp = `
helm install {{ .RepoName }}/{{ .ChartName }}-{{ .ChartVersion }}.tgz --name {{ .Name }} --namespace {{ .Namespace }} --set {{ .Values }} --debug 
`
