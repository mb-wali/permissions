package logger

import (
	"github.com/sirupsen/logrus"
)

// Log refers to the logger instance used by the permissions service.
var Log = logrus.WithFields(logrus.Fields{
	"service": "permissions",
	"art-id":  "permissions",
	"group":   "org.cyverse",
})

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
