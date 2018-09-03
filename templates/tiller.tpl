wget http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
mkdir -p /opt/bin
cp /tmp/linux-amd64/helm /opt/bin/helm
chmod +x /opt/bin/helm
/opt/bin/helm init