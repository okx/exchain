package pkg

import (
	"runtime"
	"time"

	//"github.com/satori/go.uuid"
)


func RunFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

/*
func UniqIndex() string{
	id := uuid.NewV4()
	ids := id.String()
	return ids
}*/

func GetNowTimeMs() int64{
	return time.Now().UnixNano() / 1e6
}