package expired

import (
	"os"
	"testing"
)

func Test_printHeightStatistics(t *testing.T) {
	stats := map[int64]int{
		12: 100,
		13: 1000,
		19: 10000,
		18: 200,
	}

	printHeightStatistics(stats, os.Stdout)
}
