sudo mkdir -p /root/.ssh
sudo chmod 700 /root/.ssh
sudo cd /root/.ssh
sudo touch authorized_keys
sudo chmod 600 authorized_keys
sudo bash -c "cat {{PublicKey}} >> authorized_keys"
