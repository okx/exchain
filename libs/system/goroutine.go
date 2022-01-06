package system

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type GoRoutineID int

const (
	FlagEnableGid = "enable-gid"
)

func Sleep(seconds time.Duration)  {
	time.Sleep(seconds*time.Second)
}

var (
	goroutineSpace = []byte("goroutine ")
	EnableGid      = false
)

var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

var GoRId GoRoutineID = 0

func (base GoRoutineID) String() string {
	if !EnableGid {
		return "NA"
	}
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]

	// extract the 2021 out of "goroutine 2021 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return "invalid goroutine id"
	}
	b = b[:i]

	s := string(b)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err.Error()
	}

	if int(base) == 0 {
		return s
	} else {
		return strconv.FormatUint(n, int(base))
	}
}
