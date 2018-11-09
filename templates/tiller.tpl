
sudo wget http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /opt/bin/helm
sudo chmod +x /opt/bin/helm

sudo kubectl create serviceaccount -n kube-system tiller
{{if .RBACEnabled }}
sudo kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
{{ end }}

echo "Install tiller and wait for it to be ready"
sudo /opt/bin/helm init --service-account tiller --wait
