// logger/logger.go
package logger

import (
	"os"

	"github.com/sbelyanin/SSCAN/config"
	"github.com/sirupsen/logrus"
)

func InitLogger(config config.LoggerConfig) error {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(getLogLevel(config.Level))
	logrus.SetFormatter(getFormatter(config.Format))

	return nil
}

func getLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

func getFormatter(format string) logrus.Formatter {
	switch format {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{FullTimestamp: true}
	}
}
