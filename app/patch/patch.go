package patch

import (
	"strconv"
	"sync"
)

var (
	VERSION_0_16_5_1_HEIGHT_STR       = "0"
	VERSION_0_16_5_1_HEIGHT     int64 = 0
	once                        sync.Once
)

func initVersionBlockHeight() {
	once.Do(func() {
		var err error
		if len(VERSION_0_16_5_1_HEIGHT_STR) == 0 {
			VERSION_0_16_5_1_HEIGHT_STR = "0"
		}
		VERSION_0_16_5_1_HEIGHT, err = strconv.ParseInt(VERSION_0_16_5_1_HEIGHT_STR, 10, 64)
		if err != nil {
			panic(err)
		}
	})
}

func init() {
	initVersionBlockHeight()
}