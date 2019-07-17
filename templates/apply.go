package templates

const applyTpl = `
sudo bash -c "cat <<EOF | kubectl apply -f -
{{ .Data }}
EOF"
`
