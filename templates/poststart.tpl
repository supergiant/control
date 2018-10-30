echo "PostStart started"

{{ if .IsMaster }}
    until $(curl --output /dev/null --silent --head --fail http://{{ .Host }}:{{ .Port }}); do printf '.'; sleep 5; done
    curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://{{ .Host }}:{{ .Port }}/api/v1/namespaces
    sudo kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    sudo kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    sudo kubectl config use-context default-system

    sudo kubectl create -f /etc/kubernetes/addons/kube-dns/kube-dns.yaml

    {{if .RBACEnabled }}

    sudo bash -c "cat << EOF > cluster-roles.yaml
    apiVersion: v1
    items:
    - aggregationRule:
        clusterRoleSelectors:
        - matchLabels:
            rbac.authorization.k8s.io/aggregate-to-admin: "true"
      apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: admin
      rules: null
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: cluster-admin
      rules:
      - apiGroups:
        - '*'
        resources:
        - '*'
        verbs:
        - '*'
      - nonResourceURLs:
        - '*'
        verbs:
        - '*'
    - aggregationRule:
        clusterRoleSelectors:
        - matchLabels:
            rbac.authorization.k8s.io/aggregate-to-edit: "true"
      apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
          rbac.authorization.k8s.io/aggregate-to-admin: "true"
        name: edit
      rules: null
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
          rbac.authorization.k8s.io/aggregate-to-admin: "true"
        name: system:aggregate-to-admin
      rules:
      - apiGroups:
        - authorization.k8s.io
        resources:
        - localsubjectaccessreviews
        verbs:
        - create
      - apiGroups:
        - rbac.authorization.k8s.io
        resources:
        - rolebindings
        - roles
        verbs:
        - create
        - delete
        - deletecollection
        - get
        - list
        - patch
        - update
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
          rbac.authorization.k8s.io/aggregate-to-edit: "true"
        name: system:aggregate-to-edit
      rules:
      - apiGroups:
        - ""
        resources:
        - pods/attach
        - pods/exec
        - pods/portforward
        - pods/proxy
        - secrets
        - services/proxy
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - serviceaccounts
        verbs:
        - impersonate
      - apiGroups:
        - ""
        resources:
        - pods
        - pods/attach
        - pods/exec
        - pods/portforward
        - pods/proxy
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - configmaps
        - endpoints
        - persistentvolumeclaims
        - replicationcontrollers
        - replicationcontrollers/scale
        - secrets
        - serviceaccounts
        - services
        - services/proxy
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - apps
        resources:
        - daemonsets
        - deployments
        - deployments/rollback
        - deployments/scale
        - replicasets
        - replicasets/scale
        - statefulsets
        - statefulsets/scale
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - autoscaling
        resources:
        - horizontalpodautoscalers
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - batch
        resources:
        - cronjobs
        - jobs
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - extensions
        resources:
        - daemonsets
        - deployments
        - deployments/rollback
        - deployments/scale
        - ingresses
        - networkpolicies
        - replicasets
        - replicasets/scale
        - replicationcontrollers/scale
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - policy
        resources:
        - poddisruptionbudgets
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
      - apiGroups:
        - networking.k8s.io
        resources:
        - networkpolicies
        verbs:
        - create
        - delete
        - deletecollection
        - patch
        - update
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
          rbac.authorization.k8s.io/aggregate-to-view: "true"
        name: system:aggregate-to-view
      rules:
      - apiGroups:
        - ""
        resources:
        - configmaps
        - endpoints
        - persistentvolumeclaims
        - pods
        - replicationcontrollers
        - replicationcontrollers/scale
        - serviceaccounts
        - services
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - bindings
        - events
        - limitranges
        - namespaces/status
        - pods/log
        - pods/status
        - replicationcontrollers/status
        - resourcequotas
        - resourcequotas/status
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - namespaces
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - apps
        resources:
        - daemonsets
        - deployments
        - deployments/scale
        - replicasets
        - replicasets/scale
        - statefulsets
        - statefulsets/scale
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - autoscaling
        resources:
        - horizontalpodautoscalers
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - batch
        resources:
        - cronjobs
        - jobs
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - extensions
        resources:
        - daemonsets
        - deployments
        - deployments/scale
        - ingresses
        - networkpolicies
        - replicasets
        - replicasets/scale
        - replicationcontrollers/scale
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - policy
        resources:
        - poddisruptionbudgets
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - networking.k8s.io
        resources:
        - networkpolicies
        verbs:
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:auth-delegator
      rules:
      - apiGroups:
        - authentication.k8s.io
        resources:
        - tokenreviews
        verbs:
        - create
      - apiGroups:
        - authorization.k8s.io
        resources:
        - subjectaccessreviews
        verbs:
        - create
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:aws-cloud-provider
      rules:
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
        - patch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:basic-user
      rules:
      - apiGroups:
        - authorization.k8s.io
        resources:
        - selfsubjectaccessreviews
        - selfsubjectrulesreviews
        verbs:
        - create
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:certificates.k8s.io:certificatesigningrequests:nodeclient
      rules:
      - apiGroups:
        - certificates.k8s.io
        resources:
        - certificatesigningrequests/nodeclient
        verbs:
        - create
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:certificates.k8s.io:certificatesigningrequests:selfnodeclient
      rules:
      - apiGroups:
        - certificates.k8s.io
        resources:
        - certificatesigningrequests/selfnodeclient
        verbs:
        - create
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:csi-external-attacher
      rules:
      - apiGroups:
        - ""
        resources:
        - persistentvolumes
        verbs:
        - get
        - list
        - patch
        - update
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - storage.k8s.io
        resources:
        - volumeattachments
        verbs:
        - get
        - list
        - patch
        - update
        - watch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - get
        - list
        - patch
        - update
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:csi-external-provisioner
      rules:
      - apiGroups:
        - ""
        resources:
        - persistentvolumes
        verbs:
        - create
        - delete
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims
        verbs:
        - get
        - list
        - patch
        - update
        - watch
      - apiGroups:
        - storage.k8s.io
        resources:
        - storageclasses
        verbs:
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - get
        - list
        - patch
        - update
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:discovery
      rules:
      - nonResourceURLs:
        - /api
        - /api/*
        - /apis
        - /apis/*
        - /healthz
        - /openapi
        - /openapi/*
        - /swagger-2.0.0.pb-v1
        - /swagger.json
        - /swaggerapi
        - /swaggerapi/*
        - /version
        - /version/
        verbs:
        - get
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:heapster
      rules:
      - apiGroups:
        - ""
        resources:
        - events
        - namespaces
        - nodes
        - pods
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - extensions
        resources:
        - deployments
        verbs:
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:kube-aggregator
      rules:
      - apiGroups:
        - ""
        resources:
        - endpoints
        - services
        verbs:
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:kube-controller-manager
      rules:
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - endpoints
        - secrets
        - serviceaccounts
        verbs:
        - create
      - apiGroups:
        - ""
        resources:
        - secrets
        verbs:
        - delete
      - apiGroups:
        - ""
        resources:
        - configmaps
        - endpoints
        - namespaces
        - secrets
        - serviceaccounts
        verbs:
        - get
      - apiGroups:
        - ""
        resources:
        - endpoints
        - secrets
        - serviceaccounts
        verbs:
        - update
      - apiGroups:
        - authentication.k8s.io
        resources:
        - tokenreviews
        verbs:
        - create
      - apiGroups:
        - '*'
        resources:
        - '*'
        verbs:
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:kube-dns
      rules:
      - apiGroups:
        - ""
        resources:
        - endpoints
        - services
        verbs:
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:kube-scheduler
      rules:
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - endpoints
        verbs:
        - create
      - apiGroups:
        - ""
        resourceNames:
        - kube-scheduler
        resources:
        - endpoints
        verbs:
        - delete
        - get
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - pods
        verbs:
        - delete
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - bindings
        - pods/binding
        verbs:
        - create
      - apiGroups:
        - ""
        resources:
        - pods/status
        verbs:
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - replicationcontrollers
        - services
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - apps
        - extensions
        resources:
        - replicasets
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - apps
        resources:
        - statefulsets
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - policy
        resources:
        - poddisruptionbudgets
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims
        - persistentvolumes
        verbs:
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:kubelet-api-admin
      rules:
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - proxy
      - apiGroups:
        - ""
        resources:
        - nodes/log
        - nodes/metrics
        - nodes/proxy
        - nodes/spec
        - nodes/stats
        verbs:
        - '*'
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:node
      rules:
      - apiGroups:
        - authentication.k8s.io
        resources:
        - tokenreviews
        verbs:
        - create
      - apiGroups:
        - authorization.k8s.io
        resources:
        - localsubjectaccessreviews
        - subjectaccessreviews
        verbs:
        - create
      - apiGroups:
        - ""
        resources:
        - services
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - create
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes/status
        verbs:
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - delete
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - pods
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - pods
        verbs:
        - create
        - delete
      - apiGroups:
        - ""
        resources:
        - pods/status
        verbs:
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - pods/eviction
        verbs:
        - create
      - apiGroups:
        - ""
        resources:
        - configmaps
        - secrets
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims
        - persistentvolumes
        verbs:
        - get
      - apiGroups:
        - ""
        resources:
        - endpoints
        verbs:
        - get
      - apiGroups:
        - certificates.k8s.io
        resources:
        - certificatesigningrequests
        verbs:
        - create
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims/status
        verbs:
        - get
        - patch
        - update
      - apiGroups:
        - ""
        resources:
        - serviceaccounts/token
        verbs:
        - create
      - apiGroups:
        - storage.k8s.io
        resources:
        - volumeattachments
        verbs:
        - get
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:node-bootstrapper
      rules:
      - apiGroups:
        - certificates.k8s.io
        resources:
        - certificatesigningrequests
        verbs:
        - create
        - get
        - list
        - watch
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:node-problem-detector
      rules:
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
      - apiGroups:
        - ""
        resources:
        - nodes/status
        verbs:
        - patch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:node-proxier
      rules:
      - apiGroups:
        - ""
        resources:
        - endpoints
        - services
        verbs:
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - nodes
        verbs:
        - get
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:persistent-volume-provisioner
      rules:
      - apiGroups:
        - ""
        resources:
        - persistentvolumes
        verbs:
        - create
        - delete
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims
        verbs:
        - get
        - list
        - update
        - watch
      - apiGroups:
        - storage.k8s.io
        resources:
        - storageclasses
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - watch
      - apiGroups:
        - ""
        resources:
        - events
        verbs:
        - create
        - patch
        - update
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
        name: system:volume-scheduler
      rules:
      - apiGroups:
        - ""
        resources:
        - persistentvolumes
        verbs:
        - get
        - list
        - patch
        - update
        - watch
      - apiGroups:
        - storage.k8s.io
        resources:
        - storageclasses
        verbs:
        - get
        - list
        - watch
      - apiGroups:
        - ""
        resources:
        - persistentvolumeclaims
        verbs:
        - get
        - list
        - patch
        - update
        - watch
    - aggregationRule:
        clusterRoleSelectors:
        - matchLabels:
            rbac.authorization.k8s.io/aggregate-to-view: "true"
      apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        annotations:
          rbac.authorization.kubernetes.io/autoupdate: "true"
        creationTimestamp: null
        labels:
          kubernetes.io/bootstrapping: rbac-defaults
          rbac.authorization.k8s.io/aggregate-to-edit: "true"
        name: view
      rules: null
    kind: List
    metadata: {}
    EOF"

    sudo kubectl create -f cluster-roles.yaml --validate=false
    sudo kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
    sudo kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
    sudo kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
    sudo kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
    sudo kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
    {{end}}
{{ else }}
    until $([ $(sudo docker ps |grep hyperkube| wc -l) -eq 2 ]); do printf '.'; sleep 5; done

    sudo kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
    sudo kubectl config set-context default-system --cluster=default-cluster --user=default-admin
    sudo kubectl config use-context default-system

    sudo kubectl --kubeconfig /etc/kubernetes/worker-kubeconfig.yaml config set-credentials kubelet --client-certificate /etc/kubernetes/ssl/worker.pem --client-key /etc/kubernetes/ssl/worker-key.pem --server=https://{{ .Host }}
{{ end }}

echo "PostStart finished"