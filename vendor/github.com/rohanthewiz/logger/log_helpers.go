package logger

import (
	"fmt"
	"runtime"
)

// This function is deprecated please use serr.FuncLoc()
func FunctionLoc() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}
