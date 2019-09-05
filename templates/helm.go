package templates

const helmTpl = `
echo "Installing helm"

sudo wget -nv http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /usr/bin/helm
sudo chmod +x /usr/bin/helm
sudo helm init --client-only
`
