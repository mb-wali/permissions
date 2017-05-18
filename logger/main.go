package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.WithFields(logrus.Fields{
	"service": "permissions",
	"art-id":  "permissions",
	"group":   "org.cyverse",
})

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
