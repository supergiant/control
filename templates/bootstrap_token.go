package templates

const bootstrapTokenTpl = `
{{ if .IsBootstrap }}
sudo kubeadm token create {{ .Token }} --ttl 0
# Bind uploaded certs secret to bootstrap token
sudo kubeadm init phase upload-certs --experimental-upload-certs --certificate-key {{ .CertificateKey }}
{{ end }}
`
