package logger

import (
	"fmt"
	"time"
	"strings"
	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

// Special logging for errors and structured errors (github.com/rohanthewiz/serr)
func LogErr(err error, fields ...string) {
	if err == nil {
		logrus.Error("cowardly refusing to log a nil error at ", serr.FuncLoc(2))
		return
	}
	msgs := []string{}  // for "msg" fields
	errs := []string{}  // for "error" fields

	// Add standard logging fields
	flds := logrus.Fields{"level": "error"}
	if seq, ok := flds["seq"]; !ok || seq == "" {  // set a sequence if not already set
		flds["seq"] = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if app, ok := flds["app"]; !ok || app == "" {
		flds["app"] = logOptions.AppName
	}
	if env, ok := flds["env"]; !ok || env == "" {
		flds["env"] = logOptions.Environment
	}

	// Wrap - let serr handle validation
	err = serr.LogWrap(err,  serr.CallerIndirection.GrandParent, fields...)

	// Add error string from original error
	if er := err.Error(); er != "" {
		errs = []string{er}
	}

	// If error is structured error, get key vals
	if ser, ok := err.(serr.SErr); ok {
		for key, val := range ser.FieldsMap() {
			if key != "" {
				switch strings.ToLower(key) {
				case "error":
					errs = append(errs, val)
				case "msg":
					msgs = append(msgs, val)
				default:
					flds[key] = val
				}
			}
		}
	}
	// Populate the "error" field
	if len(errs) > 0 {
		flds["error"] = strings.Join(errs, " - ")
	}
	if len(msgs) == 0 {
		msgs = []string{err.Error()}
	}
	// Log it
	logrus.WithFields(flds).Error(strings.Join(msgs, " - "))
}
