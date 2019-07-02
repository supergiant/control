package templates

const uncordonTpl = `
NODENAME=$(sudo kubectl get no -o wide|grep {{ .PrivateIP }}| awk '{ print $1 }')

if [ -z $NODENAME ]
then
	exit 0
fi

sudo kubectl uncordon $NODENAME
`
