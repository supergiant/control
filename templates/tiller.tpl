
sudo wget http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /opt/bin/helm
sudo chmod +x /opt/bin/helm
sudo /opt/bin/helm init