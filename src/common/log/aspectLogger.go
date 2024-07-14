package log

import (
	"context"
	"finance-service/common/log/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"runtime"
	"time"
)

type key string

const TraceIDKey key = "traceID"

// Function type for any function that doesn't take arguments and returns an error
type FunctionWithNoArgs func()

// LogWraps function that wraps another function
func LogWraps(fn FunctionWithNoArgs) FunctionWithNoArgs {
	funcName := getFunctionName(fn)
	return func() {
		start := time.Now()
		log.Printf("Started %s at %s", funcName, start)
		fn()
		end := time.Now()
		log.Printf("Finished %s at %s, duration: %s", funcName, end, end.Sub(start))
	}
}

func ExtractTraceIDFromContextOrUnknown(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return "unknown"
	}
	return traceID
}

// GetTraceIDFromContext extracts the trace Id from the gin.Context.
func GetTraceIDFromGinContextOrUnknown(c *gin.Context) string {
	if traceID, exists := c.Get(string(TraceIDKey)); exists {
		if strTraceID, ok := traceID.(string); ok {
			return strTraceID
		}
	}
	return "unknown"
}

// getFunctionName returns the name of a function
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// TraceIDMiddleware adds a trace Id to the context and logs the request start and end
func TraceIDMiddleware(logger *CommonLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(TraceIDKey), traceID)

		start := time.Now()
		logger.Log(dto.LogEntry{
			Level:   logrus.InfoLevel,
			TraceId: traceID,
			Event:   "request_started",
			Message: "Request started",
			Context: map[string]interface{}{
				"method": c.Request.Method,
				"url":    c.Request.URL.Path,
			},
		})

		c.Next()

		end := time.Now()
		logger.Log(dto.LogEntry{
			Level:   logrus.InfoLevel,
			TraceId: traceID,
			Event:   "request_finished",
			Message: "Request finished",
			Context: map[string]interface{}{
				"method":          c.Request.Method,
				"url":             c.Request.URL.Path,
				"requestDuration": end.Sub(start),
			},
		})
	}
}

// LogFunctionExecution logs the start and end of a function execution, only for DEBUG purposes
func LogFunctionExecution(logger *CommonLogger, ctx context.Context, fn func() error) error {
	fnName := getFunctionName(fn)

	traceID := ctx.Value(TraceIDKey).(string)
	start := time.Now()
	logger.Log(dto.LogEntry{
		Level:   logrus.DebugLevel,
		TraceId: traceID,
		Event:   "function_start",
		Message: fnName + " started",
	})

	err := fn()

	end := time.Now()
	logger.Log(dto.LogEntry{
		Level:   logrus.DebugLevel,
		TraceId: traceID,
		Event:   "function_start",
		Message: fnName + " started",
		Context: map[string]interface{}{
			"duration": end.Sub(start),
		},
	})

	return err
}
