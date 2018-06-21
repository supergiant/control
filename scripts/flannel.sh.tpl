#!/bin/bash
echo "[Unit]
Description=Networking service
Requires=etcd-member.service
[Service]
Environment=FLANNEL_IMAGE_TAG=v0.9.0
ExecStartPre=/usr/bin/etcdctl set /coreos.com/network/config '{"Network":"10.2.0.0/16", "Backend": {"Type": "vxlan"}}'" > \
/etc/systemd/system/flannel.service
systemctl start flannel
