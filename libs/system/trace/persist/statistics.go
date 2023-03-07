package persist

import (
	"fmt"

	"github.com/okx/okbchain/libs/system/trace"
)

var stats *statistics

func init() {
	stats = &statistics{
		BaseStatistics: trace.NewSummary(),
	}
}

type statistics struct {
	trace.BaseStatistics
}

func GetStatistics() *statistics {
	return stats
}

func (s *statistics) Format() string {
	var res string
	for _, tag := range s.GetTags() {
		res += fmt.Sprintf("%s<%dms>, ", tag, s.GetValue(tag)/1e6)
	}

	return res[0 : len(res)-2]
}
