## Logging package

### Usage

```go
package main
import (
    "errors"
    . "github.com/rohanthewiz/logger
    "github.com/rohanthewiz/serr
)

func main() {
	InitLog(LogOptions{
		AppName: "My App",
		Environment: "development",
		Format: "text",
		Level: "info",
		InfoPath: "log/info.log",
		ErrorPath: "log/error.log",
	})
	defer CloseLog()

    Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
    Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2")

	err := errors.New("This is the original error")

	// We can log a standard error, the message will be err.Error()
	LogErr(err, "message")
		//=> ERRO[0000] message	error="This is the original error"

	// Multiple arguments after message are treated as a key, value list and will wrap the error 
	// Be careful to use pairs of fields after message.
	LogErr(err, "message", "key1", "value1", "key2", "value2")
		//=> ERRO[0000] message error="This is the original error" key1=value1 key2=value2

	// See log_test.go for more examples	
}
```
