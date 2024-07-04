package dto

import "github.com/sirupsen/logrus"

type LogEntry struct {
	Level      logrus.Level
	TraceID    string
	Event      string
	Message    string
	Context    map[string]interface{}
	StackTrace error
}
