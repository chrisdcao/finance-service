package log

import (
	"finance-service/utils"
	"finance-service/utils/log/dto"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type CommonLogger struct {
	serviceName string
	logger      *logrus.Logger
}

func NewCommonLogger() *CommonLogger {
	serviceName := utils.GetPackagePath()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout) // TODO: Ensure logs go to stdout for Logstash to pick up

	return &CommonLogger{
		serviceName: serviceName,
		logger:      logger,
	}
}

func (this *CommonLogger) Log(entry dto.LogEntry) {
	logEntry := this.logger.WithFields(logrus.Fields{
		"timestamp":  time.Now().Format(time.RFC3339),
		"level":      entry.Level.String(),
		"service":    this.serviceName,
		"traceID":    entry.TraceID,
		"event":      entry.Event,
		"message":    entry.Message,
		"context":    entry.Context,
		"stackTrace": fmt.Sprintf("%+v", entry.StackTrace), // Include the stack trace
	})
	logEntry.Log(entry.Level, entry.Message)
}
