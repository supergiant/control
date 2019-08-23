package templates

const addAuthorizedKeysTpl = `
{{ if .UserName }}
sudo adduser {{ .UserName }} --gecos "{{ .UserName }},{{ .UserName }},{{ .UserName }},{{ .UserName }}" --disabled-password

sudo mkdir -p /home/{{ .UserName }}/.ssh
sudo chmod 700 /home/{{ .UserName }}/.ssh
sudo touch /home/{{ .UserName }}/.ssh/authorized_keys
sudo chmod 600 /home/{{ .UserName }}/.ssh/authorized_keys

sudo chown -R {{ .UserName }} /home/{{ .UserName }}/.ssh/
sudo chown {{ .UserName }} /home/{{ .UserName }}/.ssh/authorized_keys

sudo bash -c "cat << EOF >> /home/{{ .UserName }}/.ssh/authorized_keys
{{ .PublicKey }}
EOF"

echo "{{ .UserName }} ALL=(ALL:ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/{{ .UserName }}
{{ end }}

sudo mkdir -p /root/.ssh
sudo chmod 700 /root/.ssh
sudo touch /root/.ssh/authorized_keys
sudo chmod 600 /root/.ssh/authorized_keys

sudo bash -c "cat << EOF >> /root/.ssh/authorized_keys
{{ .PublicKey }}
EOF"
`
