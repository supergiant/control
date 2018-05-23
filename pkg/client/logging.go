package client

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	logger := logrus.New()
	Log = logger
}
