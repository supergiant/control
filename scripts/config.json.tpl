{
  "bind-address": "{{ .Kube.MasterPrivateIP }}",
  "hostname-override": "{{ .Kube.MasterPrivateIP }}",
  "cluster-cidr": "172.30.0.0/16",
  "logtostderr": true,
  "v": 0,
  "allow-privileged": true,
  "master": "http://{{ .Kube.MasterPrivateIP }}:8080",
  "etcd-servers": "http://{{ .Kube.MasterPrivateIP }}:2379"
}
