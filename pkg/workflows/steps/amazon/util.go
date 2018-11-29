package amazon

import (
	"context"
	"time"
	"net/http"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
)

func FindOutboundIP(ctx context.Context, findExternalIP func() (string, error)) (string, error) {
	var (
		timeout  = time.Second * 10
		attempts = 5
		publicIP string
		err      error
	)

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
	serviceURL := "http://myexternalip.com/raw"

	resp, err := http.Get(serviceURL)
	if err != nil {
		logrus.Debugf("error while accessing %s", serviceURL)
		return "", err
	}

	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	publicIP := string(res)
	publicIP = strings.Replace(publicIP, "\n", "", -1)

	return publicIP, err
}

