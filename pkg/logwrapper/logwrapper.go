package logwrapper

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{}
	standardLogger.SetReportCaller(true)
	standardLogger.SetLevel(logrus.DebugLevel)

	file, err := os.OpenFile("shodanmonitor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		standardLogger.Fatal(err)
		return nil
	}

	mw := io.MultiWriter(os.Stdout, file)
	standardLogger.SetOutput(mw)
	return standardLogger
}
