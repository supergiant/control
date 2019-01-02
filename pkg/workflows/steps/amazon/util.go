package amazon

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	timeout  = time.Second * 10
	serviceURLs = []string{
		"http://checkip.amazonaws.com/",
		"http://myexternalip.com/raw",
	}

	attempts = 5
	publicIP string

)

func FindOutboundIP(ctx context.Context, findExternalIP func() (string, error)) (string, error) {
	var err error

	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			publicIP, err = findExternalIP()

			if err == nil {
				return publicIP, nil
			} else {
				logrus.Debugf("attempt #%d Error getting public IP sleep for %v",
					i+1, timeout)
				time.Sleep(timeout)
				timeout = timeout * 2
			}
		}
	}

	return publicIP, err
}

func findOutBoundIP() (string, error) {
	var (
		publicIP string
		resp 	*http.Response
		err      error
	)

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 10,
	}
	for _, serviceURL := range serviceURLs {
		resp, err = client.Get(serviceURL)
		if err != nil {
			logrus.Debugf("error while accessing %s", serviceURL)
			continue
		}

		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		publicIP = string(res)
		publicIP = strings.Replace(publicIP, "\n", "", -1)

		logrus.Debugf("get IP from %s", serviceURL)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		break
	}

	return publicIP, err
}
