// Copyright (c) 2012-2013 Matt Nunogawa @amattn
// This source code is release under the MIT License, http://opensource.org/licenses/MIT

package deeperror

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var globalErrorLoggingEnabled bool

func init() {
	globalErrorLoggingEnabled = false
}

const (
	globalDefaultStatusCode = http.StatusInternalServerError
)

//
type DeepError struct {
	Num           int64
	Filename      string
	CallingMethod string
	Line          int
	EndUserMsg    string
	DebugMsg      string
	DebugFields   map[string]interface{}
	Err           error // inner or source error
	StatusCode    int
	StackTrace    string
}

// Primary Constructor.  Create a DeepError ptr with the given number, end user message and optional parent error.
func New(num int64, endUserMsg string, parentErr error) *DeepError {
	const skipStackFrames = 2
	return NewS(num, endUserMsg, parentErr, skipStackFrames)
}

// Wrapper version.  Create a DeepError with the starting stack frame higher in the stack with
// ptr with the given number, end user message and optional parent error. A value of 1 for skipStackFrames
// would yield the callers frame
func NewS(num int64, endUserMsg string, parentErr error, skipStackFrames int) *DeepError {
	e := new(DeepError)
	e.Num = num
	e.EndUserMsg = endUserMsg
	e.Err = parentErr
	e.StatusCode = globalDefaultStatusCode
	e.DebugFields = make(map[string]interface{})

	gerr, ok := parentErr.(*DeepError)
	if ok {
		if gerr != nil {
			e.StatusCode = gerr.StatusCode
		}
	}

	pc, file, line, ok := runtime.Caller(skipStackFrames)

	if ok {
		e.Line = line
		components := strings.Split(file, "/")
		e.Filename = components[(len(components) - 1)]
		f := runtime.FuncForPC(pc)
		e.CallingMethod = f.Name()
	}

	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)

	e.StackTrace = string(buf[:n])

	if globalErrorLoggingEnabled {
		log.Print(e)
	}
	return e
}

// HTTP variant.  Create a DeepError with the given http status code
func NewHTTPError(num int64, endUserMsg string, err error, statusCode int) *DeepError {
	derr := New(num, endUserMsg, err)
	derr.StatusCode = statusCode
	if len(endUserMsg) == 0 {
		derr.EndUserMsg = http.StatusText(statusCode)
	}
	return derr
}

// Convenience method.  creates a simple DeepError with the given error number.  The error message is set to "TODO"
func NewTODOError(num int64, printArgs ...interface{}) *DeepError {
	derr := New(num, "TODO", nil)

	for i, printArg := range printArgs {
		derr.AddDebugField(strconv.Itoa(i), printArg)
	}

	return derr
}

// Convenience method.  This will return nil if parrentErr == nil.  Otherwise it will create a DeepError and return that.
func NewOrNilFromParent(num int64, endUserMsg string, parentErr error) error {
	if parentErr == nil {
		return nil
	}
	return New(num, endUserMsg, parentErr)
}

// Convenience method.  Equivalient to derr:=New(...); log.Fatal(derr)
func Fatal(num int64, endUserMsg string, parentErr error) {
	derr := New(num, endUserMsg, parentErr)
	log.Fatal(derr)
}

// Add arbitrary debugging data to a given DeepError
func (derr *DeepError) AddDebugField(key string, value interface{}) {
	derr.DebugFields[key] = value
}

// cConform to the new Unwrap interface.
// Unwrap() will expose errors further down the error chain
// This should allow support for Is() and As() in Go 1.13 and later
// Alternatively, earlier version of Go can
// import "golang.org/x/xerrors" to get library support
// for Is(), As(), and Unwrap()
// see https://blog.golang.org/go1.13-errors for details
func (derr *DeepError) Unwrap() error {
	return derr.Err
}

// internal usage for formatting/pretty printing
func prependToLines(para, prefix string) string {
	lines := strings.Split(para, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// Check if the current status code matches the global default
func (derr *DeepError) StatusCodeIsDefaultValue() bool {
	if derr.StatusCode == globalDefaultStatusCode {
		return true
	} else {
		return false
	}
}

// Conform to the go built-in error interface
// http://golang.org/pkg/builtin/#error
func (derr *DeepError) Error() string {

	parentError := "nil"

	if derr.Err != nil {
		parentError = prependToLines(derr.Err.Error(), "-- ")
	}

	debugFieldStrings := make([]string, 0, len(derr.DebugFields))
	for k, v := range derr.DebugFields {
		str := fmt.Sprintf("\n-- DebugField[%s]: %+v", k, v)
		debugFieldStrings = append(debugFieldStrings, str)
	}

	dbgMsg := ""
	if len(derr.DebugMsg) > 0 {
		dbgMsg = "\n-- DebugMsg: " + derr.DebugMsg
	}

	return fmt.Sprintln(
		"\n\n-- DeepError",
		derr.Num,
		derr.StatusCode,
		derr.Filename,
		derr.CallingMethod,
		"line:", derr.Line,
		"\n-- EndUserMsg: ", derr.EndUserMsg,
		dbgMsg,
		strings.Join(debugFieldStrings, ""),
		"\n-- StackTrace:",
		strings.TrimLeft(prependToLines(derr.StackTrace, "-- "), " "),
		"\n-- ParentError:", parentError,
	)
}

// enable/disable automatic logging of deeperrors upon creation
func ErrorLoggingEnabled() bool {
	return globalErrorLoggingEnabled
}

// anything performed in this anonymous function will not trigger automatic logging of deeperrors upon creation
type NoErrorsLoggingAction func()

// you can use this method to temporarily disable automatic logging of deeperrors
func ExecWithoutErrorLogging(action NoErrorsLoggingAction) {
	// this is racy...  I feel ashamed.
	original := globalErrorLoggingEnabled
	globalErrorLoggingEnabled = false
	action()
	globalErrorLoggingEnabled = original
}
