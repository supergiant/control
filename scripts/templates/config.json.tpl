{
  "bind-address": "{{ .MasterPrivateIP }}",
  "hostname-override": "{{ .MasterPrivateIP }}",
  "cluster-cidr": "172.30.0.0/16",
  "logtostderr": true,
  "v": 0,
  "allow-privileged": true,
  "master": "http://{{ .MasterPrivateIP }}:8080",
  "etcd-servers": "http://{{ .MasterPrivateIP }}:2379"
}
