package utils

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)


type GoroutineID int

var goroutineSpace = []byte("goroutine ")

var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

var GoId GoroutineID = 0

func (base GoroutineID) String() string {
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Extract the `1606` out of "goroutine 1606 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return "N/A"
	}
	b = b[:i]

	// check if it is valid
	s := string(b)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return "WRFMT"
	}

	// 0 mean "used as it like"
	if int(base) == 0 {
		return s
	} else {
		return strconv.FormatUint(n, int(base))
	}
}
