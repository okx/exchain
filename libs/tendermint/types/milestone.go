package types

import "strconv"

var (
	VENUS_HEIGHT string
	venusHeight  int64
)

func init() {
	if VENUS_HEIGHT == "" {
		return
	}
	n, err := strconv.ParseInt(VENUS_HEIGHT, 10, 64)
	if err != nil {
		panic(err)
	}
	venusHeight = n
}

func HigherThanVenus(height int64) bool {
	if venusHeight == 0 {
		return false
	}
	return height > venusHeight
}
