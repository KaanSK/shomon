package logwrapper

import (
	"os"

	"github.com/sirupsen/logrus"
)

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// Logger : Globally shared logging instance
var Logger = NewLogger()

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.SetLevel(logrus.InfoLevel)

	standardLogger.SetOutput(os.Stdout)
	return standardLogger
}
