package templates

const tillerTpl = `
echo "Installing tiller and waiting for it to be ready"

sudo kubectl create serviceaccount -n kube-system tiller
{{if .RBACEnabled }}
sudo kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
{{ end }}

sudo /usr/bin/helm init --automount-service-account-token --wait
`
