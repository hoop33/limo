package logger

import (
	"sync"
	"time"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"github.com/rohanthewiz/rotatelogs"
)

type LogOptions struct {
	AppName     string
	Environment string
	Format      string
	Level       string
	InfoPath    string
	ErrorPath   string
}

const logsChannelSize = 2000
const errsChannelSize = 200
var logsChannel chan [][]byte  // receive [] of []byte
var errsChannel chan []string
var logsDone chan bool
var logsWaitGroup = new(sync.WaitGroup) // this will serve both error and nonerror logs
var logOptions LogOptions

// This is our main entry point for logging
func InitLog(lopts LogOptions) {
	// Init some package variables
	logOptions = lopts  // make options available to the logging package in logOptions
	logsChannel = make(chan [][]byte, logsChannelSize)
	errsChannel = make(chan []string, errsChannelSize)
	logsDone = make(chan bool)

	initLogrus()

	// Start the log listener
	go pollForLogs(logsDone)
}

// Close out asynchronous logging
func CloseLog() {
	logsWaitGroup.Wait()  // fan in all log goroutines
	// Close the channel so no more logs can be sent and the log poller knows to start wrapping up
	close(logsChannel)
	close(errsChannel)
	<- logsDone // wait for the poller to completely wrap up
}

func initLogrus() {
	// Set Formatter
	if strings.ToLower(logOptions.Format) == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// Create log rotater

	// Info
	infoLogPath, err := filepath.Abs(logOptions.InfoPath)
	if err != nil {
		fmt.Println("Info log path invalid -", err); os.Exit(1)
	}
	infoWriter, err := rotatelogs.New(infoLogPath + ".%Y%m%d%H%M",
			rotatelogs.WithLinkName(infoLogPath),
			rotatelogs.WithMaxAge(time.Duration(4*24) * time.Hour),  // Keep for 4 days
			rotatelogs.WithRotationTime(time.Duration(2*24) * time.Hour),  // rotate every
	)
	if err != nil {
		fmt.Println("failed to create info log file:", err); os.Exit(1)
	}

	// Error
	errorLogPath, err := filepath.Abs(logOptions.ErrorPath)
	if err != nil {
		fmt.Println("Error log path invalid -", err); os.Exit(1)
	}
	errWriter, err := rotatelogs.New(errorLogPath + ".%Y%m%d%H%M",
			rotatelogs.WithLinkName(errorLogPath),
			rotatelogs.WithMaxAge(time.Duration(14*24) * time.Hour),
			rotatelogs.WithRotationTime(time.Duration(4*24) * time.Hour),
	)
	if err != nil {
		fmt.Println("failed to create error log file:", err); os.Exit(1)
	}

	// Add LFS hook
	logrus.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: infoWriter,
		logrus.InfoLevel: infoWriter,
		logrus.WarnLevel: errWriter,
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
	}))

	switch strings.ToLower(logOptions.Level) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}