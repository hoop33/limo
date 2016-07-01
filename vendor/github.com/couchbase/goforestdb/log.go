package forestdb

//#include "log.h"
import "C"
import (
	"log"
	"unsafe"
)

//export LogCallbackInternal
func LogCallbackInternal(errCode C.int, msg *C.char, ctx *C.char) {
	context := (*C.log_context)(unsafe.Pointer(ctx))
	cb := logCallbacks[context.offset]
	userCtx := logContexts[context.offset]
	cb(C.GoString(context.name), int(errCode), C.GoString(msg), userCtx)
}

//export FatalErrorCallbackInternal
func FatalErrorCallbackInternal() {
	fatalErrorCallback()
}

// Logger interface
type Logger interface {
	// Warnings, logged by default.
	Warnf(format string, v ...interface{})
	// Errors, logged by default.
	Errorf(format string, v ...interface{})
	// Fatal errors. Will not terminate execution.
	Fatalf(format string, v ...interface{})
	// Informational messages.
	Infof(format string, v ...interface{})
	// Timing utility
	Debugf(format string, v ...interface{})
	// Program execution tracing. Not logged by default
	Tracef(format string, v ...interface{})
}

type Dummy struct {
}

func (*Dummy) Fatalf(_ string, _ ...interface{}) {
}

func (*Dummy) Errorf(_ string, _ ...interface{}) {
}

func (*Dummy) Warnf(_ string, _ ...interface{}) {
}

func (*Dummy) Infof(_ string, _ ...interface{}) {
}

func (*Dummy) Debugf(_ string, _ ...interface{}) {
}

func (*Dummy) Tracef(_ string, _ ...interface{}) {
}

type LogLevel int

const (
	LogFatal LogLevel = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

type LeveledLog struct {
	level LogLevel
}

func NewLeveledLog(level LogLevel) *LeveledLog {
	return &LeveledLog{level: level}
}

func (l *LeveledLog) Fatalf(format string, a ...interface{}) {
	if l.level >= LogFatal {
		log.Fatalf(format, a...)
	}
}

func (l *LeveledLog) Errorf(format string, a ...interface{}) {
	if l.level >= LogError {
		log.Printf(format, a...)
	}
}

func (l *LeveledLog) Warnf(format string, a ...interface{}) {
	if l.level >= LogWarn {
		log.Printf(format, a...)
	}
}

func (l *LeveledLog) Infof(format string, a ...interface{}) {
	if l.level >= LogInfo {
		log.Printf(format, a...)
	}
}

func (l *LeveledLog) Debugf(format string, a ...interface{}) {
	if l.level >= LogDebug {
		log.Printf(format, a...)
	}
}

func (l *LeveledLog) Tracef(format string, a ...interface{}) {
	if l.level >= LogTrace {
		log.Printf(format, a...)
	}
}

// Logger to use
var Log Logger = &Dummy{}
