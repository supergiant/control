sudo kubectl drain $(sudo kubectl get no -o wide|grep {{ .PrivateIP }}| awk '{ print $1 }') \
--ignore-daemonsets --force --delete-local-data

