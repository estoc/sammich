package chewcrew

import (
	"encoding/json"
	"fmt"
	logging "github.com/op/go-logging"
	"runtime"
)

// Globally available server log.
var serverLog *logger = nil

const (
	logFormat = "%{color}%{level} %{module} %{color:reset}%{message}"
)

// A simple wrapper around the go-logging library that allows us to easily create child loggers.
//
// TODO: support filesystem logging
type logger struct {
	libLogger *logging.Logger
	level     int
}

// Get a new server logger.
func NewServerLogger(logLevel int, logFormat string) *logger {
	libLogger := logging.MustGetLogger("server")
	logger := &logger{libLogger, logLevel}

	logging.SetFormatter(logging.MustStringFormatter(logFormat))

	logging.SetLevel(logging.Level(logLevel), "server")
	return logger
}

// Get a child logger.
func (l *logger) Child(name string) (*logger, error) {
	libLogger, err := logging.GetLogger(name)
	if err != nil {
		return nil, err
	}

	child := &logger{libLogger, l.level}
	logging.SetLevel(logging.Level(l.level), name)
	return child, nil
}

// Log at critical level. See *logger.logAtLevel for usage details
func (l *logger) Critical(args ...interface{}) {
	l.logAtLevel("CRITICAL", getLogCallerPath(2), args...)
}

// Log at error level. See *logger.logAtLevel for usage details
func (l *logger) Error(args ...interface{}) {
	l.logAtLevel("ERROR", getLogCallerPath(2), args...)
}

// Log at warning level. See *logger.logAtLevel for usage details
func (l *logger) Warning(args ...interface{}) {
	l.logAtLevel("WARNING", getLogCallerPath(2), args...)
}

// Log at notice level. See *logger.logAtLevel for usage details
func (l *logger) Notice(args ...interface{}) {
	l.logAtLevel("NOTICE", getLogCallerPath(2), args...)
}

// Log at info level. See *logger.logAtLevel for usage details
func (l *logger) Info(args ...interface{}) {
	l.logAtLevel("INFO", getLogCallerPath(2), args...)
}

// Log at debug level. See *logger.logAtLevel for usage details
func (l *logger) Debug(args ...interface{}) {
	l.logAtLevel("DEBUG", getLogCallerPath(2), args...)
}

// Internally used method that allows us to leverage the go-logging library while enabling a more
// robust and expressive logging interface. Each logging level (see public methods above) has three
// available interfaces:
//   - logAtLevel(format string, args interface{}...)
//     - log a format string. args is optional.
//   - logAtLevel(err error)
//     - log an error. the error's message. if using main.Error, a stack and context message (if set)
//     will also be provided.
//   - logAtLevel(struct interface{}, format string, args interface{}...)
//     - log a struct alongside a format string. args is optional. the struct will be marshalled to
//     json so that it can be pretty printed. if marshalling fails, an error will be logged at the
//     server level and will immediately return.
func (l *logger) logAtLevel(level string, caller string, args ...interface{}) {
	var format string
	var formatArgs []interface{}

	switch args[0].(type) {
	case string: // recvd formatted string [and args]
		format = "%s: " + args[0].(string)
		formatArgs = append([]interface{}{caller}, args[1:]...)
	case error: // recvd an error
		format = "%s: %s"
		formatArgs = []interface{}{caller, args[0].(error).Error()}
	default: // recvd a struct, formatted string [and args].. will attempt to marshall to json
		marshalled, err := json.MarshalIndent(args[0], "", "    ")
		if err != nil {
			serverLog.Error("an error occurred when marshalling a struct for logging: %s", err)
			return
		}
		format = "%s: " + args[1].(string) + "\n%s"
		formatArgs = append(append([]interface{}{caller}, args[2:]...), marshalled)
	}

	switch level {
	case "CRITICAL":
		l.libLogger.Critical(format, formatArgs...)
	case "ERROR":
		l.libLogger.Error(format, formatArgs...)
	case "WARNING":
		l.libLogger.Warning(format, formatArgs...)
	case "NOTICE":
		l.libLogger.Notice(format, formatArgs...)
	case "INFO":
		l.libLogger.Info(format, formatArgs...)
	case "DEBUG":
		l.libLogger.Debug(format, formatArgs...)
	}
}

// Allows us to retrieve our caller's context
func getLogCallerPath(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if ok != true {
		file = "log caller not recoverable"
		line = -1
	}

	return fmt.Sprintf("%s:%v", file, line)
}
