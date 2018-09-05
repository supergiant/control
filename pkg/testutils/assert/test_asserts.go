package assert

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	//DefaultETCDURL used in integration tests
	DefaultETCDURL = "http://127.0.0.1:2379"
)

//etcdRunning is test assert function useful in integration tests to assert that test dependencies are running
func MustRunETCD(etcdHost string) {
	if err := CheckETCD(etcdHost); err != nil {
		logrus.Fatal(err)
	}
}

func CheckETCD(etcdHost string) error {
	logrus.Infof("connecting to the ETCD...")
	cl, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdHost},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer cl.Close()
	_, err = cl.Dial(etcdHost)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}
