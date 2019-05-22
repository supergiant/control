package templates

const addAuthorizedKeysTpl = `
sudo mkdir -p /root/.ssh
sudo chmod 700 /root/.ssh
sudo touch /root/.ssh/authorized_keys
sudo chmod 600 /root/.ssh/authorized_keys

sudo bash -c "cat << EOF >> /root/.ssh/authorized_keys
{{ .PublicKey }}
EOF"
`
