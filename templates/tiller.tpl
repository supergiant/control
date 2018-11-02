
sudo wget http://storage.googleapis.com/kubernetes-helm/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/
sudo tar -C /tmp -xvf /tmp/helm-v{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo cp /tmp/linux-amd64/helm /opt/bin/helm
sudo chmod +x /opt/bin/helm

sudo kubectl create serviceaccount -n kube-system tiller
sudo kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
sudo /opt/bin/helm init --service-account tiller

until $([ $(sudo kubectl get pods --namespace=kube-system|grep tiller|grep Running|wc -l) -eq 1 ]); do printf '.'; sleep 5; done
TILLER_PORT=44135
TILLER_POD=$(sudo kubectl get pods --namespace=kube-system|grep tiller|awk '{print $1}')
TILLER_IP=$(sudo kubectl describe pod $TILLER_POD -n kube-system|grep IP| awk '{print $2}')
until $(curl --output /dev/null --silent --head --fail http://$TILLER_IP:$TILLER_PORT/readiness); do
    printf '.'
    sleep 5
done
