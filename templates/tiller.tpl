echo "Installing tiller and waiting for it to be ready"

sudo wget -nv http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /usr/bin/helm
sudo chmod +x /usr/bin/helm

sudo kubectl create serviceaccount -n kube-system tiller
{{if .RBACEnabled }}
sudo kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
{{ end }}

sudo /usr/bin/helm init --automount-service-account-token --wait