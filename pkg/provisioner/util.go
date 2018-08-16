package provisioner

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type TokenGetter interface {
	GetToken(int) (string, error)
}

type EtcdTokenGetter struct {
	discoveryUrl string
}

func NewEtcdTokenGetter() *EtcdTokenGetter {
	return &EtcdTokenGetter{
		discoveryUrl: "https://discovery.etcd.io/new?size=%d",
	}
}

func (e *EtcdTokenGetter) GetToken(num int) (string, error) {
	resp, err := http.Get(fmt.Sprintf(e.discoveryUrl, num))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
