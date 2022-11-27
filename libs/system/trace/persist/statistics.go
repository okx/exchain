package persist

import (
	"fmt"
	"sync"

	"github.com/okex/exchain/libs/system/trace"
)

var once sync.Once
var stats *statistics

type statistics struct {
	trace.BaseStatistics
}

func GetStatistics() *statistics {
	once.Do(func() {
		stats = &statistics{
			BaseStatistics: trace.NewSummary(),
		}
	})

	return stats
}

func (s *statistics) Format() string {
	var res string
	for _, tag := range s.GetTags() {
		res += fmt.Sprintf("%s<%dms>, ", tag, s.GetValue(tag)/1e6)
	}

	return res[0 : len(res)-2]
}
