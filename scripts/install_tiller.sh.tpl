wget http://storage.googleapis.com/kubernetes-helm/{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
tar -C /tmp -xvf /tmp/{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
cp /tmp/linux-amd64/helm /opt/bin/helm
chmod +x /opt/bin/helm
/opt/bin/helm init