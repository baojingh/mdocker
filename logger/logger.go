package logger

import (
	logger "github.com/sirupsen/logrus"
	rotate  "github.com/lestrrat-go/file-rotatelogs"
	"path/filepath"
)

var log = logger.New()

func init() {
	path := "./"
	rotate.New(
		filepath.Join()


	)


}

func New() *logger.Logger() {
	return log
}



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
