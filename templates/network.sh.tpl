sudo curl -L {{ .EtcdRepositoryUrl }}/v{{ .EtcdVersion }}/etcd-v{{ .EtcdVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz -o /tmp/etcd-v{{ .EtcdVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz
sudo tar xzvf /tmp/etcd-v{{ .EtcdVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz -C /usr/bin --strip-components=1

ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_INITIAL_ADVERTISE_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCDCTL_API=3 /usr/bin/etcdctl version

sudo /usr/bin/etcdctl set /coreos.com/network/config '{"Network":"{{ .Network }}", "Backend": {"Type": "{{ .NetworkType }}"}}'
sudo /usr/bin/etcdctl get /coreos.com/network/config
