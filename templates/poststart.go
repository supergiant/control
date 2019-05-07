package templates

const poststartTpl = `
{{ if .IsBootstrap }}
{{ if .RBACEnabled }}
sudo kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
sudo kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
sudo kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
sudo kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
{{ end }}
{{ end }}

sudo bash -c "cat << EOF > /etc/security/limits.conf
root soft  nofile 300000
root hard  nofile 300000
EOF"

if [[ $(whoami) != root ]]; then
  sudo cp -r /home/$(whoami)/.kube /root/
fi
`
