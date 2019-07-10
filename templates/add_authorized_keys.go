package templates

const addAuthorizedKeysTpl = `
sudo bash -c "cat << EOF >> /home/ubuntu/.ssh/authorized_keys
{{ .PublicKey }}
EOF"

sudo adduser {{ .UserName }} --gecos "{{ .UserName }},{{ .UserName }},{{ .UserName }},{{ .UserName }}" --disabled-password
sudo bash -c "cat > /etc/sudoers <<EOF
{{ .UserName }} ALL=(ALL) NOPASSWD:ALL
EOF"

sudo mkdir -p /home/{{ .UserName }}/.ssh
sudo chmod 700 /home/{{ .UserName }}/.ssh
sudo touch /home/{{ .UserName }}/.ssh/authorized_keys
sudo chmod 600 /home/{{ .UserName }}/.ssh/authorized_keys

sudo chown -R {{ .UserName }} /home/{{ .UserName }}/.ssh/
sudo chown {{ .UserName }} /home/{{ .UserName }}/.ssh/authorized_keys

sudo bash -c "cat << EOF >> /home/{{ .UserName }}/.ssh/authorized_keys
{{ .PublicKey }}
EOF"

sudo bash -c "cat << EOF >> /home/{{ .UserName }}/.ssh/authorized_keys
{{ .BootstrapPublicKey }}
EOF"
`
