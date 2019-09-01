package logger

import (
	"errors"
	"github.com/rohanthewiz/serr"
	"testing"
)

// This is just a visual test for now (commented outputs are abbreviated)
func TestLogging(t *testing.T) {
	InitLog(LogOptions{
		AppName: "Log Test",
		Environment: "Dev",
		Level: "Debug",
	})
	defer CloseLog()

	t.Log("Dummy use of testing")

	err := errors.New("This is the original error")

	// We can log a standard error, the message will be err.Error()
	LogErr(err, "message")
	//=> ERRO[0000] message error="This is the original error" level=error location="logger/logger_test.go:16"

	// Single argument after err becomes part of logrus message
	LogErr(err, "Custom message here", "error", "I'm making this up")
	//=> ERRO[0000] Custom message here error="This is the original error - I'm making this up" level=error location="logger/logger_test.go:20"
	// Multiple arguments after err are treated as a key, value list and will wrap the error
	LogErr(err, "message", "key1", "value1", "key2", "value2")
	//=> ERRO[0000] msg="message" error="This is the original error" key1=value1 key2=value2 level=error location="logger/logger_test.go:23"

	// Multiple arguments after err are treated as a key, value list and will wrap the error
	LogErr(err, "This is an error", "error", "Error Code: ABCDE321",
		"msg", "This is a critical error", "key1", "value1")
	//=> ERRO[0000] This is an error - This is a critical error   app= env= error="This is the original error - Error Code: ABCDE321" key1=value1 level=error location="logger/logger_test.go:28"

	err2 := serr.Wrap(err, "Gosh! We got an error!")
	LogErr(err2)
	//=> ERRO[0000] message - Gosh! We got an error!              app= env= error="This is the original error" level=error location="logger/logger_test.go:32 - logger/logger_test.go:31"

	// We can log an SErr wrapped error
	err3 := serr.Wrap(err2, "cat", "aight", "dogs", "I dunno")
	LogErr(err3, "Animals, do we really need them? Yes!!", "author", "me")
	//=> ERRO[0000] Animals, do we really need them? Yes!! - Gosh! We got an error!  app= author=me cat=aight dogs="I dunno" env= error="This is the original error" level=error location="logger/logger_test.go:37 - logger/logger_test.go:36 - logger/logger_test.go:31"
	LogErrAsync(err2, "We can now log errors asynchronously!", "key1", "value1", "key2", "value")
	LogAsync("Warn", "Let see how this warning goes", "keyA", "valueA")
}
