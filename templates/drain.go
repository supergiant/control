package templates

const drainTpl = `
NODENAME=$(sudo kubectl get no -o wide|grep {{ .PrivateIP }}| awk '{ print $1 }')

if [ -z "NODENAME" ]
then
	exit 0
fi

sudo kubectl drain $NODENAME \
--ignore-daemonsets --force --delete-local-data

sudo kubectl delete no $NODENAME
`
