package expired

import (
	"os"
	"testing"
)

func Test_printStatistics(t *testing.T) {
	stats := map[string]int{
		"0xabc": 100,
		"0xbcd": 1000,
		"0xecd": 10000,
		"0xead": 200,
	}

	printStatistics(stats, os.Stdout)
}
