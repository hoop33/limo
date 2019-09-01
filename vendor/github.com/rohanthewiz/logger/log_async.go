package logger

import (
	"fmt"
	"time"
	"errors"
)

// Queue up logs on logsChannel. Level, message, and all arguments are strings
func LogAsync(level string, msg string, args... string) {
	LogAsyncBin(level, msg, nil, args...)
}

// The transport here allows for a binary third argument, all other arguments are strings
func LogAsyncBin(level string, msg string, bin *[]byte, args ...string) {
	blank_bytes := []byte{}
	if bin == nil {
		bin = &blank_bytes
	}
	argsSlice := [][]byte{ []byte(level), []byte(msg), *bin } // first section of arguments
	for _, arg := range args {  // get the rest of the arguments
		argsSlice = append(argsSlice, []byte(arg))
	}
	// Lock in a sequence attribute here before the async call
	argsSlice = append(argsSlice, []byte("seq"))
	argsSlice = append(argsSlice, []byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	logsWaitGroup.Add(1)  // track the number of log senders
	go func(slice [][]byte) {
		logsChannel <- slice // send to the channel. // A go routine is needed bc the logsChannel may be blocked (full)
		logsWaitGroup.Done()     // one less log sender
	}(argsSlice)
}

// Poll the LogsChannel for incoming messages of [][]byte
// Arguments are the receive only logs channel and send only done channel
func pollForLogs(done chan <- bool) {
	defer func() {  // Flush any hooks here
		done <- true  // signal caller when we are done
	}()
	var logsComplete, errsComplete bool  // The channel's processing is complete
	for {
		select {  // Select can multiplex cases reading from multiple channels . We will block till there is a message
					// on a channel: one of the cases unblocks
		case attrs, ok := <-logsChannel: //
			if !ok { // the channel is closed *and* empty, so wrap up
				if errsComplete { return }
				logsComplete = true
			} else if len(attrs) >= 3 {
				LogBinary(string(attrs[0]), string(attrs[1]), attrs[2], attrs[3:]...) // receive the item and call Log()
			}
		case errAttrs, ok := <- errsChannel:
			if !ok {
				if logsComplete { return }
				errsComplete = true
			} else if len(errAttrs) > 0 {
				LogErr(errors.New(errAttrs[0]), errAttrs[1:]...)
			}

		// we don't timeout. Logs run for the life of the app
		}
	}
}

func LogErrAsync(err error, args ... string) {
	argsSlice := []string{err.Error()}
	argsSlice = append(argsSlice, args...)
	logsWaitGroup.Add(1)
	go func(slice []string) {
		errsChannel <- slice
		logsWaitGroup.Done()
	}(argsSlice)
}