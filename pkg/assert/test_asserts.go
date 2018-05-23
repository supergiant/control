package assert

import (
	"time"

	"github.com/coreos/etcd/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//etcdRunning is test assert function useful in integration tests to assert that test dependencies are running
func EtcdRunning(etcdHost string) {
	cl, err := client.New(client.Config{
		Endpoints: []string{etcdHost},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err != nil {
		logrus.Fatal(err)
	}
	defer cancel()
	_, err = cl.GetVersion(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
}
