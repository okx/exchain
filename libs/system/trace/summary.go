package trace

import (
	"fmt"
	systemlog "log"

	//"github.com/okex/exchain/libs/tendermint/libs/log"
)

var sum *Summary
type Summary struct {
	statisticMap map[string]int64
	keys []string
}

func insertElapse(tag string, elapse int64)  {
	if sum == nil {
		return
	}
	sum.insert(tag, elapse)
}

func GetTraceSummary() *Summary {
	once.Do(func() {
		sum = &Summary{
			statisticMap: make(map[string]int64),
		}
	})
	return sum
}

func (s *Summary) Init(keys ...string)  {
	for _, k := range keys {
		s.statisticMap[k] = 0
	}
	s.keys = keys
}

//func (s *Summary) Dump(logger log.Logger)  {
//	for _, k := range s.keys {
//		logger.With("module", "main").Info("Summary", k, s.statisticMap[k])
//	}
//}

func (s *Summary) Dump(context string)  {
	var res string
	for _, k := range s.keys {
		res += fmt.Sprintf("%s=%d, ", k, s.statisticMap[k])
	}
	systemlog.Println("Elapse Summary", context, res)
}

func (s *Summary) insert(tag string, elapse int64)  {
	_, ok := s.statisticMap[tag]
	if !ok {
		return
	}
	s.statisticMap[tag] += elapse
}