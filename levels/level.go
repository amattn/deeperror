package levels

// The deeperror convention is: TRACE < DEBUG < INFO < WARN < ERROR < FATAL

// Levels are ints
// Values are strings

type ErrorLevel int
type ErrorLevelString string

const (
	// Undefined is the default value for unintialized variables.
	Undefined ErrorLevel = iota // 0
	Off                         // 1
	Fatal                       // 2 ...
	Error
	Warn
	Info
	Debug
	Trace
	// we reserver the right to be more loquacious in future versions
)

const (
	UNDEF ErrorLevelString = "UNDEF"
	OFF   ErrorLevelString = "OFF"
	FATAL ErrorLevelString = "FATAL"
	ERROR ErrorLevelString = "ERROR"
	WARN  ErrorLevelString = "WARN"
	INFO  ErrorLevelString = "INFO"
	DEBUG ErrorLevelString = "DEBUG"
	TRACE ErrorLevelString = "TRACE"

	// only for when things have gone terribly wrong:
	UNKNOWN = "<UNKNOWN ERROR LEVEL STRING>"
)

func (el ErrorLevel) String() string {
	els := ErrorLevelToErrorLevelString(el)
	return els.String()
}

func (els ErrorLevelString) String() string {
	return string(els)
}

func ErrorLevelStringToLevel(els ErrorLevelString) ErrorLevel {
	switch els {
	case UNDEF:
		return Undefined
	case OFF:
		return Off
	case FATAL:
		return Fatal
	case ERROR:
		return Error
	case WARN:
		return Warn
	case INFO:
		return Info
	case DEBUG:
		return Debug
	case TRACE:
		return Trace
	default:
		return -1
	}
}

func ErrorLevelToErrorLevelString(el ErrorLevel) ErrorLevelString {
	switch el {
	case Undefined:
		return UNDEF
	case Off:
		return OFF
	case Fatal:
		return FATAL
	case Error:
		return ERROR
	case Warn:
		return WARN
	case Info:
		return INFO
	case Debug:
		return DEBUG
	case Trace:
		return TRACE
	default:
		return UNKNOWN
	}
}
