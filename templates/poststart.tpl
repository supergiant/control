echo "PostStart started"

{{ if .IsMaster }}
    until $(sudo curl --output /dev/null --silent --head --fail http://{{ .Host }}:{{ .Port }}); do printf '.'; sleep 5; done
    sudo curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://{{ .Host }}:{{ .Port }}/api/v1/namespaces
    sudo kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    sudo kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    sudo kubectl config use-context default-system

    kubectl create -f /etc/kubernetes/addons/kube-dns/kube-dns.yaml

    {{if .RBACEnabled }}
    sudo kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
    sudo kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
    sudo kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
    sudo kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
    sudo kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
    {{end}}
{{ else }}
    until $([ $(sudo docker ps |grep hyperkube| wc -l) -eq 2 ]); do printf '.'; sleep 5; done

    sudo kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    sudo kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    sudo kubectl config use-context default-system

    sudo kubectl --kubeconfig /etc/kubernetes/worker-kubeconfig.yaml config set-credentials kubelet --client-certificate /etc/kubernetes/ssl/worker.pem --client-key /etc/kubernetes/ssl/worker-key.pem --server=https://{{ .Host }}
    sudo bash -c "until $([ $(sudo docker ps |grep hyperkube| wc -l) -eq 2 ]); do printf '.'; sleep 5; done"
{{ end }}

echo "PostStart finished"