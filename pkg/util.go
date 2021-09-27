package pkg

import (
	"runtime"
	"strings"
	"time"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	names := strings.Split(f.Name(), "/")
	return names[len(names) - 1]
}

func GetNowTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetNowTimeNs() int64 {
	return time.Now().UnixNano()
}
