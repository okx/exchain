package pkg

import (
	"runtime"
	"time"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func GetNowTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}
