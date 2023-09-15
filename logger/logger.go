package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	logger "github.com/sirupsen/logrus"
)

var log = logger.New()

func init() {
	path := "./"
	writer, err := rotate.New(
		filepath.Join(path, fmt.Sprintf("mdocker-%s.log", "%Y%m%d")),
		rotate.WithLinkName(filepath.Join(path, "mdocker.log")),
		rotate.WithRotationCount(5),
		rotate.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.SetOutput(io.MultiWriter(writer))
	}
}

func New() *logger.Logger {
	return log
}
