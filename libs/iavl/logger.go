package iavl

const (
	FlagOutputModules = "iavl-output-modules"
)

const (
	IavlErr   = 0
	IavlInfo  = 1
	IavlDebug = 2
)

var (
	logFunc LogFuncType = nil

	OutputModules map[string]int

	iavlLogger logger
)

type logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

type LogFuncType func(level int, format string, args ...interface{})

func SetLogFunc(l LogFuncType) {
	logFunc = l
}

func SetLogger(l logger) {
	iavlLogger = l
}

func iavlLog(module string, level int, msg string, keyvals ...interface{}) {
	if v, ok := OutputModules[module]; ok && v != 0 && iavlLogger != nil {
		switch level {
		case IavlErr:
			iavlLogger.Error(msg, keyvals...)
		case IavlInfo:
			iavlLogger.Info(msg, keyvals...)
		case IavlDebug:
			iavlLogger.Debug(msg, keyvals...)
		}
	}
}
