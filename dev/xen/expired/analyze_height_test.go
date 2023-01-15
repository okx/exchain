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
	mintBlockTime := map[int64]int64{
		12: 16788888,
		13: 16788888,
		19: 16788888,
		18: 16788888,
	}
	claim := map[int64]int{
		17: 5,
	}
	claimBlockTime := map[int64]int64{
		17: 16788888,
	}

	printHeightStatistics(stats, mintBlockTime, claim, claimBlockTime, os.Stdout)
}
