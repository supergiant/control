# install supergiant/capacity chart

KUBESCALER_CONFIG=/tmp/kubescaler.conf
USERDATA_FILE=/tmp/userdata.txt

cat <<CLUSTERSERVICESEOF | sudo tee $KUBESCALER_CONFIG
{{ .KubescalerConfig }}
CLUSTERSERVICESEOF

# replace every new line character with \\n
cat <<CLUSTERSERVICESEOF | sed ':a;N;$!ba;s/\n/\\n/g' | sudo tee $USERDATA_FILE
{{ .Userdata }}
CLUSTERSERVICESEOF

sudo /opt/bin/helm repo add supergiant https://supergiant.github.io/charts

sudo /opt/bin/helm install supergiant/capacity \
    --version {{ .Version }} \
    --name=capacity \
    --namespace=kube-system \
    --set-file config.kubescaler.raw=$KUBESCALER_CONFIG \
    --set-file config.userdata=$USERDATA_FILE