echo '[Unit]
      Description=Kubernetes Kubelet Server
      Documentation=https://github.com/kubernetes/kubernetes
      Requires=docker.service network-online.target
      After=docker.service network-online.target

      [Service]
      ExecStartPre=/bin/mkdir -p /var/lib/kubelet
      ExecStartPre=/bin/mount --bind /var/lib/kubelet /var/lib/kubelet
      ExecStartPre=/bin/mount --make-shared /var/lib/kubelet
      ExecStart=/usr/bin/docker run \
              --net=host \
              --pid=host \
              --privileged \
              -v /dev:/dev \
              -v /sys:/sys:ro \
              -v /var/run:/var/run:rw \
              -v /var/lib/docker/:/var/lib/docker:rw \
              -v /var/lib/kubelet/:/var/lib/kubelet:shared \
              -v /var/log:/var/log:shared \
              -v /srv/kubernetes:/srv/kubernetes:ro \
              -v /etc/kubernetes:/etc/kubernetes:ro \
              gcr.io/google-containers/hyperkube:v{{ .KubernetesVersion }} \
              /hyperkube kubelet --allow-privileged=true \
              --cluster-dns=10.3.0.10 \
              --cluster_domain=cluster.local \
              --pod-manifest-path=/etc/kubernetes/manifests \
              --kubeconfig=/etc/kubernetes/worker-kubeconfig.yaml \
              --volume-plugin-dir=/etc/kubernetes/volumeplugins \
              {{- .KubeProviderString }}
              --fail-swap-on=false \
              --register-node=true
      Restart=always
      StartLimitInterval=0
      RestartSec=10
      KillMode=process

      [Install]
      WantedBy=multi-user.target' > /etc/systemd/system/kubelet.service
systemctl start kubelet