#!/bin/bash
wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld
chmod 755 /usr/bin/flanneld
sed -i 's/REPLACEME/'`ifconfig|grep "10\."|grep "inet "|cut -f10 -d" "`'/g' /etc/default/flanneld

cat << EOF > /etc/systemd/system/flanneld.service
[Unit]
Description=Networking service
Requires=etcd-member.service
[Service]
Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
ExecStartPre=/usr/bin/etcdctl set /coreos.com/network/config '{"Network":"{{ .Network }}", "Backend": {"Type": "{{ .NetworkType }}"}}'
EOF
systemctl enable flanneld.service
systemctl restart flanneld.service
