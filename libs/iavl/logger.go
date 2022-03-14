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
)

type LogFuncType func(level int, format string, args ...interface{})

func SetLogFunc(l LogFuncType) {
	logFunc = l
}

func iavlLog(module string, level int, format string, args ...interface{}) {
	if v, ok := OutputModules[module]; ok && v != 0 && logFunc != nil {
		logFunc(level, format, args...)
	}
}
