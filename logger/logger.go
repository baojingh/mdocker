package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	Log     *logrus.Logger
	initLog sync.Once
)

func init() {
	initLog.Do(func() {
		Log = logrus.New()
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:           "2006-01-02 15:04:05",
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
			DisableLevelTruncation:    true,
		})
	})
}
