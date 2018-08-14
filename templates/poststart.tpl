echo "PostStart started"
until $(curl --output /dev/null --silent --head --fail http://{{ .Host }}:{{ .Port }}); do printf '.'; sleep 5; done
curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://{{ .Host }}:{{ .Port }}/api/v1/namespaces
/opt/bin/kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
/opt/bin/kubectl config set-context default-system --cluster=default-cluster --user=default-admin
/opt/bin/kubectl config use-context default-system

{{if .RBACEnabled }}
/opt/bin/kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
/opt/bin/kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
/opt/bin/kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
/opt/bin/kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
/opt/bin/kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
{{end}}

echo "PostStart finished"