
sudo wget http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /opt/bin/helm
sudo chmod +x /opt/bin/helm

sudo kubectl create serviceaccount -n kube-system tiller
sudo kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
sudo /opt/bin/helm init --service-account tiller