echo "PostStart started"


sudo bash -c "cat << EOF > /etc/security/limits.conf
root soft  nofile 300000
root hard  nofile 300000
EOF"

if [[ $(whoami) != root ]]; then
  sudo cp -r /home/$(whoami)/.kube /root/
fi

echo "PostStart finished"