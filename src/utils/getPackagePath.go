package utils

import (
	"runtime"
	"strings"
)

// GetPackagePath returns the full package path of the caller.
func GetPackagePath() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()
	// Extract the package path from the function name
	lastSlash := strings.LastIndex(funcName, "/")
	if lastSlash == -1 {
		return funcName
	}
	return funcName[:lastSlash]
}
