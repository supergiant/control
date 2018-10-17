echo "PostStart started"

{{ if .IsMaster }}
    until $(curl --output /dev/null --silent --head --fail http://{{ .Host }}:{{ .Port }}); do printf '.'; sleep 5; done
    curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://{{ .Host }}:{{ .Port }}/api/v1/namespaces
    kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    kubectl config use-context default-system

    kubectl create -f /etc/kubernetes/addons/kube-dns/kube-dns.yaml

    {{if .RBACEnabled }}
    kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
    kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
    kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
    kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
    kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
    {{end}}
{{ else }}
    until $([ $(docker ps |grep hyperkube| wc -l) -eq 2 ]); do printf '.'; sleep 5; done

    kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    kubectl config use-context default-system

    kubectl --kubeconfig /etc/kubernetes/worker-kubeconfig.yaml config set-credentials kubelet --client-certificate /etc/kubernetes/ssl/worker.pem --client-key /etc/kubernetes/ssl/worker-key.pem --server=https://{{ .Host }}
{{ end }}

echo "PostStart finished"