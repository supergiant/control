package templates

const bootstrapTokenTpl = `
{{ if .IsBootstrap }}
sudo kubeadm token create {{ .Token }} --ttl 0
# Bind uploaded certs secret to bootstrap token

{{ if not .IsImport }}
	sudo kubeadm init phase upload-certs --upload-certs --certificate-key {{ .CertificateKey }}
{{ end }}
{{ end }}
`
