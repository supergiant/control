package provisioner

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EtcdTokenGetter struct {
	discoveryUrl string
}

func NewEtcdTokenGetter() *EtcdTokenGetter {
	return &EtcdTokenGetter{
		discoveryUrl: "https://discovery.etcd.io/new?size=%d",
	}
}

func (e *EtcdTokenGetter) GetToken(ctx context.Context, num int) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(e.discoveryUrl, num), nil)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	return string(body), nil
}

func makeFileName(taskID string) string {
	return fmt.Sprintf("%s.log", taskID)
}
