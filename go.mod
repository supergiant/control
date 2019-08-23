module github.com/supergiant/control

go 1.12

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190606204050-af9c91bd2759
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190606210616-f848dc7be4a4
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190606205144-71ebb8303503
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190606212257-347f17c60af0
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.3
)

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.4 // indirect
	github.com/Azure/azure-sdk-for-go v31.1.0+incompatible
	github.com/Azure/go-autorest v11.4.0+incompatible
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.16.0+incompatible // indirect
	github.com/aokoli/goutils v1.0.1 // indirect
	github.com/apparentlymart/go-cidr v0.0.0-20180915144716-1755c023625e
	github.com/aws/aws-sdk-go v1.19.5
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/dgrijalva/jwt-go v3.1.0+incompatible
	github.com/digitalocean/godo v1.10.0
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/spdystream v0.0.0-20170912183627-bc6354cbbc29 // indirect
	github.com/elazarl/goproxy v0.0.0-20170405201442-c4fc26588b6e // indirect
	github.com/etcd-io/bbolt v0.0.0-20190125183005-8693da9f4d97
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20160516000752-02826c3e7903 // indirect
	github.com/golang/mock v1.2.0 // indirect
	github.com/golang/protobuf v1.3.0
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/google/uuid v1.1.0 // indirect
	github.com/googleapis/gnostic v0.0.0-20180419025854-664776a3b48a // indirect
	github.com/gophercloud/gophercloud v0.3.0
	github.com/gorilla/handlers v0.0.0-20170602151648-a4d79d4487c2
	github.com/gorilla/mux v1.6.2
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.6.2 // indirect
	github.com/hpcloud/tail v1.0.0
	github.com/huandu/xstrings v0.0.0-20180906151751-8bbcf2f9ccb5 // indirect
	github.com/imdario/mergo v0.3.7
	github.com/jarcoal/httpmock v0.0.0-20181110092731-53def6cd0f87
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v0.0.0-20180701071628-ab8a2e0c74be
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/pkg/errors v0.8.0
	github.com/rakyll/statik v0.1.6
	github.com/sirupsen/logrus v1.2.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/supergiant/capacity v0.0.0-20190513092134-fa714465dd86
	github.com/technosophos/moniker v0.0.0-20180509230615-a5dbd03a2245
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190621222207-cc06ce4a13d4
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/oauth2 v0.0.0-20190319182350-c85d3e98c914
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/api v0.3.0
	google.golang.org/genproto v0.0.0-20190321212433-e79c0c59cdb5 // indirect
	gopkg.in/asaskevich/govalidator.v8 v8.0.0-20171111151018-521b25f4b05f
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190703205437-39734b2a72fe
	k8s.io/apiextensions-apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/apimachinery v0.0.0-20190703205208-4cfb76a8bf76
	k8s.io/apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cloud-provider v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/cluster-bootstrap v0.0.0-20190703212826-5ad085674a4f // indirect
	k8s.io/helm v2.14.3+incompatible
	k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30 // indirect
	k8s.io/kube-proxy v0.0.0-20190703212322-69d540a3479c // indirect
	k8s.io/kubelet v0.0.0-20190704010802-f16c4cee528c // indirect
	k8s.io/kubernetes v0.0.0-00010101000000-000000000000
	sigs.k8s.io/yaml v1.1.0 // indirect
)
