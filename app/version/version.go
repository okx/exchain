package version

import (
	"strconv"
	"sync"
)

var (
	VERSION_0_16_6_1_HEIGHT               = "0"
	VERSION_0_16_6_1_HEIGHT_NUM     int64 = 0
	once                        sync.Once
)

func strin2number(input string) int64 {
	if len(input) == 0 {
		return 0
	}

	res, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func init() {
	once.Do(func() {
		VERSION_0_16_6_1_HEIGHT_NUM = strin2number(VERSION_0_16_6_1_HEIGHT)
	})
}

