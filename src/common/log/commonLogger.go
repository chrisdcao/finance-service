package log

import (
	"finance-service/common/log/dto"
	"finance-service/utils"
	"fmt"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"net"
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

	// connect to logstash hook
	logStashHookAddress, err := utils.GetEnvVariableFromKey("LOGSTASH_HOOK_ADDRESS")
	if err != nil {
		logger.Fatalf("Failed to rertieve LOGSTASH_HOOK_ADDRESS from .env: %v", err)
	}
	conn, err := net.Dial("tcp", logStashHookAddress)
	if err != nil {
		logger.Fatalf("Failed to connect to Logstash: %v", err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "myappName"}))
	logger.Hooks.Add(hook)

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout) // set the std out for logstash to get
	//logger.SetLevel(logrus.DebugLevel) // Set the appropriate log level

	return &CommonLogger{
		serviceName: serviceName,
		logger:      logger,
	}
}

func (this *CommonLogger) Log(entry dto.LogEntry) {
	logEntry := this.logger.WithFields(logrus.Fields{
		"timestamp":  time.Now().Format(time.RFC3339),
		"level":      entry.Level.String(),
		"service":    utils.GetPackagePath(),
		"traceID":    entry.TraceId,
		"event":      entry.Event,
		"message":    entry.Message,
		"context":    entry.Context,
		"stackTrace": fmt.Sprintf("%+v", entry.StackTrace), // Include the stack trace
	})
	logEntry.Log(entry.Level, entry.Message)
}
