package assert

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

//etcdRunning is test assert function useful in integration tests to assert that test dependencies are running
func EtcdRunning(etcdHost string) {
	cl, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer cl.Close()
	_, err = cl.Dial(etcdHost)
	if err != nil {
		logrus.Fatal(err)
	}
}
