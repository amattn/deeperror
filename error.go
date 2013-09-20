package deeperror

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

var gERROR_LOGGING_ENABLED bool

func init() {
	gERROR_LOGGING_ENABLED = false
}

type DeepError struct {
	Num           int64
	Filename      string
	CallingMethod string
	Line          int
	EndUserMsg    string
	DebugMsg      string
	Err           error // inner or source error
	StatusCode    int
	StackTrace    string
}

func New(num int64, endUserMsg string, parentErr error) *DeepError {
	e := new(DeepError)
	e.Num = num
	e.EndUserMsg = endUserMsg
	e.Err = parentErr
	e.StatusCode = http.StatusInternalServerError // default status code...

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
	return grunwayErrorPtr
}

func prependToLines(para, prefix string) string {
	lines := strings.Split(para, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func (e *DeepError) Error() string {

	parentError := "nil"

	// fmt.Println("THISERR", e.Num, "PARENT ERR", e.Err)

	if e.Err != nil {
		parentError = prependToLines(e.Err.Error(), "-- ")
	}

	return fmt.Sprintln(
		"\n\n-- DeepError",
		e.Num,
		e.StatusCode,
		e.Filename,
		e.CallingMethod,
		"line:", e.Line,
		"\n-- EndUserMsg: ", e.EndUserMsg,
		"\n-- DebugMsg: ", e.DebugMsg,
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
