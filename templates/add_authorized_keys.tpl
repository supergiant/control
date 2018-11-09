sudo mkdir -p /root/.ssh
sudo chmod 700 /root/.ssh
sudo touch /root/.ssh/authorized_keys
sudo chmod 600 /root/.ssh/authorized_keys
sudo bash -c "echo \"{{ .PublicKey }}\" >> /root/.ssh/authorized_keys "
