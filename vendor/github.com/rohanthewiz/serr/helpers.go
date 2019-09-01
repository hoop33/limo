package serr

import (
	"fmt"
	"runtime"
	"strings"
	"os"
)

//const (
//	Caller             = 1
//	CallersParent      = 2
//	CallersGrandParent = 3
//)
var CallerIndirection = callerIndirection{1, 2, 3, 4}
type callerIndirection struct {
	Caller, Parent, GrandParent, GreatGrandParent int
}

// Return caller or ancestors calling location
// Optional caller indirection can be
// 1 - immediate caller (default)
// 2 - the callers parent
// 3 - the callers grandparent -- you get the idea
func FuncLoc(callerIndir ...int) string {
	lvl := CallerIndirection.Caller
	if len(callerIndir) > 0 {
		lvl = callerIndir[0]
	}

	_, file, line, ok := runtime.Caller(lvl)
	if !ok {
		return "could not determine location"
	}
	return PathLevel(fmt.Sprintf("%s:%d", file, line))
}

// Return a portion of the fullpath
// 0 - file only  // e.g. main.go
// 1 - file and parent  // e.g. myproject/main.go
// 2 - file up to grandparent  // e.g. githubusername/myproject/main.go
// defaults to 1 - parent/file
func PathLevel(path string, level ...uint) (subpath string) {
	lvl := 1
	if len(level) > 0 {
		lvl = int(level[0])
	}
	if path == "" {
		return path
	}
	// On Windows filepath.Separator may not be that returned by runtime.Caller()
	// This is the case under GitBash at least so leaning towards '/'.
	sepr := '\\'
	if os.IsPathSeparator('/') {  // If any O/S says this is legit, use it
		sepr = '/'
	}
	separator := fmt.Sprintf("%c", sepr)
	tokens := strings.Split(path, separator)
	ln := len(tokens)
	if ln <= 1 {
		//fmt.Println("path not split")
		return path
	}

	idx := len(tokens) - int(lvl) - 1
	if idx < 0 {
		idx = 0
	}
	return strings.Join(tokens[idx:], separator)
}

// This function is deprecated. Please use serr.FuncLoc()
//func FunctionLoc() string {
//	_, file, line, ok := runtime.Caller(2)
//	if !ok {
//		return ""
//	}
//	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
//}
