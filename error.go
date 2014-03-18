// Copyright (c) 2012-2013 Matt Nunogawa @amattn
// This source code is release under the MIT License, http://opensource.org/licenses/MIT

package deeperror

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

var gERROR_LOGGING_ENABLED bool

func init() {
	gERROR_LOGGING_ENABLED = false
}

const (
	gDEFAULT_STATUS_CODE = http.StatusInternalServerError
)

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

func Fatal(num int64, endUserMsg string, parentErr error) {
	derr := New(num, endUserMsg, parentErr)
	log.Fatal(derr)
}

func New(num int64, endUserMsg string, parentErr error) *DeepError {
	e := new(DeepError)
	e.Num = num
	e.EndUserMsg = endUserMsg
	e.Err = parentErr
	e.StatusCode = gDEFAULT_STATUS_CODE
	e.DebugFields = make(map[string]interface{})

	gerr, ok := parentErr.(*DeepError)
	if ok {
		if gerr != nil {
			e.StatusCode = gerr.StatusCode
		}
	}

	pc, file, line, ok := runtime.Caller(1)

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

	if gERROR_LOGGING_ENABLED {
		log.Print(e)
	}
	return e
}

func NewHTTPError(num int64, endUserMsg string, err error, statusCode int) *DeepError {
	grunwayErrorPtr := New(num, endUserMsg, err)
	grunwayErrorPtr.StatusCode = statusCode
	if len(endUserMsg) == 0 {
		grunwayErrorPtr.EndUserMsg = http.StatusText(statusCode)
	}
	return grunwayErrorPtr
}

func NewTODOError(num int64) *DeepError {
	grunwayErrorPtr := New(num, "TODO", nil)
	return grunwayErrorPtr
}

func (derr *DeepError) AddDebugField(key string, value interface{}) {
	derr.DebugFields[key] = value
}

func prependToLines(para, prefix string) string {
	lines := strings.Split(para, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func (e *DeepError) StatusCodeIsDefaultValue() bool {
	if e.StatusCode == gDEFAULT_STATUS_CODE {
		return true
	} else {
		return false
	}
}

func (e *DeepError) Error() string {

	parentError := "nil"

	// fmt.Println("THISERR", e.Num, "PARENT ERR", e.Err)

	if e.Err != nil {
		parentError = prependToLines(e.Err.Error(), "-- ")
	}

	debugFieldStrings := make([]string, 0, len(e.DebugFields))
	for k, v := range e.DebugFields {
		str := fmt.Sprintf("\n-- DebugField[%s]: %+v", k, v)
		debugFieldStrings = append(debugFieldStrings, str)
	}

	dbgMsg := ""
	if len(e.DebugMsg) > 0 {
		dbgMsg = "\n-- DebugMsg: " + e.DebugMsg
	}

	return fmt.Sprintln(
		"\n\n-- DeepError",
		e.Num,
		e.StatusCode,
		e.Filename,
		e.CallingMethod,
		"line:", e.Line,
		"\n-- EndUserMsg: ", e.EndUserMsg,
		dbgMsg,
		strings.Join(debugFieldStrings, ""),
		"\n-- StackTrace:",
		strings.TrimLeft(prependToLines(e.StackTrace, "-- "), " "),
		"\n-- ParentError:", parentError,
	)
}

func ErrorLoggingEnabled() bool {
	return gERROR_LOGGING_ENABLED
}

type NoErrorsLoggingAction func()

// you can use this method to temporarily disable
func ExecWithoutErrorLogging(action NoErrorsLoggingAction) {
	// this is racy...  I feel ashamed.
	original := gERROR_LOGGING_ENABLED
	gERROR_LOGGING_ENABLED = false
	action()
	gERROR_LOGGING_ENABLED = original
}
