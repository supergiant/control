sudo set -n
sudo bash -c "cat > script <<EOF
{{ .Data }}
EOF"

sudo kubectl create capacity --from-file=script -n {{ .Namespace }}
