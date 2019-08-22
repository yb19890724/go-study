package main

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	entry := logrus.WithFields(logrus.Fields{
		"name": "test",
	})
	entry.Info("message1")
	entry.Info("message2")
	// time="2019-01-24T19:04:51+08:00" level=info msg=message1 name=test
	// time="2019-01-24T19:04:51+08:00" level=info msg=message2 name=test
}
