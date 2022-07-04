package itlogs

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Log(LogLevel, string)
	Trace(string)
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

var defaultLogger Logger = NewDefaultConsoleLogger(Debug)

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

type LogLevel int

func (l LogLevel) String() string {
	return []string{
		"TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL",
	}[l]
}

const (
	Trace LogLevel = iota
	Debug
	Info
	Warning
	Error
	Fatal
)

type LogMsgFlag int

func (f LogMsgFlag) String() string {
	return map[LogMsgFlag]string{
		Ldate:         "Ldate",
		Ltime:         "Ltime",
		Lmicroseconds: "Lmicroseconds",
		Llongfile:     "Llongfile",
		Lshortfile:    "Lshortfile",
		LUTC:          "LUTC",
		Lmsgprefix:    "Lmsgprefix",
		LstdFlags:     "LstdFlags",
	}[f]
}

const (
	Ldate         LogMsgFlag = log.Ldate
	Ltime         LogMsgFlag = log.Ltime
	Lmicroseconds LogMsgFlag = log.Lmicroseconds
	Llongfile     LogMsgFlag = log.Llongfile
	Lshortfile    LogMsgFlag = log.Lshortfile
	LUTC          LogMsgFlag = log.LUTC
	Lmsgprefix    LogMsgFlag = log.Lmsgprefix
	LstdFlags     LogMsgFlag = log.LstdFlags
)

type DefaultLogger struct {
	Level LogLevel
	trace *log.Logger
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

func NewLogger(out io.Writer, flags []LogMsgFlag, level LogLevel) *DefaultLogger {
	flag := 0

	if len(flags) == 0 {
		flag = int(Ldate | Ltime | Lmsgprefix)
	} else {
		for _, f := range flags {
			flag = flag | int(f)
		}
	}

	logger := &DefaultLogger{
		Level: level,
		trace: log.New(out, Trace.String()+" ", flag),
		debug: log.New(out, Debug.String()+" ", flag),
		info:  log.New(out, Info.String()+" ", flag),
		warn:  log.New(out, Warning.String()+" ", flag),
		err:   log.New(out, Error.String()+" ", flag),
		fatal: log.New(out, Fatal.String()+" ", flag),
	}

	return logger
}

func NewConsoleLogger(flags []LogMsgFlag, level LogLevel) *DefaultLogger {
	return NewLogger(os.Stdout, flags, level)
}

func NewDefaultConsoleLogger(level LogLevel) *DefaultLogger {
	return NewConsoleLogger([]LogMsgFlag{}, level)
}

func NewDefaultLogger(out io.Writer, level LogLevel) *DefaultLogger {
	return NewLogger(out, []LogMsgFlag{}, level)
}

func (l DefaultLogger) Log(level LogLevel, msg string) {
	if level < l.Level {
		return
	}

	withFile := func(l *log.Logger) bool {
		return l.Flags()&(int(Llongfile)|int(Lshortfile)) != 0
	}

	var logger *log.Logger

	switch level {
	case Trace:
		logger = l.trace
	case Debug:
		logger = l.debug
	case Info:
		logger = l.info
	case Warning:
		logger = l.warn
	case Error:
		logger = l.err
	case Fatal:
		logger = l.fatal
	}

	if withFile(logger) {
		logger.Output(2, msg)
	} else {
		logger.Print(msg)
	}
}

func (l DefaultLogger) Trace(msg string) {
	l.Log(Trace, msg)
}

func (l DefaultLogger) Debug(msg string) {
	l.Log(Debug, msg)
}

func (l DefaultLogger) Info(msg string) {
	l.Log(Info, msg)
}

func (l DefaultLogger) Warn(msg string) {
	l.Log(Warning, msg)
}

func (l DefaultLogger) Error(msg string) {
	l.Log(Error, msg)
}

func (l DefaultLogger) Fatal(msg string) {
	l.Log(Fatal, msg)
}
