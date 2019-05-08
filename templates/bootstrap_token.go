package templates

const bootstrapTokenTpl = `
sudo kubeadm token create {{ .Token }} --ttl 0
`
