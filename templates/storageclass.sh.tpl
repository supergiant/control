

{{ if eq .Provider "aws"}}
sudo bash -c "cat > storageclass.yaml <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: default
  labels:
    k8s-addon: storage-aws.addons.k8s.io
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2

---

apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: gp2
  annotations:
    storageclass.beta.kubernetes.io/is-default-class: \"true\"
  labels:
    k8s-addon: storage-aws.addons.k8s.io
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
EOF"
{{else if eq .Provider "gce"}}
 sudo bash -c "cat > storageclass.yaml <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: default
  annotations:
    storageclass.beta.kubernetes.io/is-default-class: "true"
  labels:
    kubernetes.io/cluster-service: "true"
    k8s-addon: storage-gce.addons.k8s.io
    addonmanager.kubernetes.io/mode: EnsureExists
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-standard
EOF"
{{ else }}
 sudo bash -c "cat > storageclass.yaml <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
EOF"
{{ end }}
echo applying default storage class
sudo cat ./storageclass.yaml
sudo kubectl apply -f storageclass.yaml


