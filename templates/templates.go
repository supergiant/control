package templates

var Default = map[string]string{
	"add_authorized_keys":        addAuthorizedKeysTpl,
	"bootstrap_token":            bootstrapTokenTpl,
	"certificates":               certificatesTpl,
	"clustercheck":               clustercheckTpl,
	"cni":                        cniTpl,
	"docker":                     dockerTpl,
	"download_kubernetes_binary": downloadKubernetesBinaryTpl,
	"drain":                      drainTpl,
	"kubeadm":                    kubeadmTpl,
	"kubelet":                    kubelet,
	"network":                    networkTpl,
	"poststart":                  poststartTpl,
	"prometheus":                 prometheusTpl,
	"storageclass":               storageclassTpl,
	"tiller":                     tillerTpl,
	"upgrade":                    upgradeTpl,
	"evacutate":                  evacuateTpl,
	"uncordon":                   uncordonTpl,
}
