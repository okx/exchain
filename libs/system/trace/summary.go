package trace

import (
	"fmt"
	//"github.com/okx/okbchain/libs/tendermint/libs/log"
)

var sum *Summary

func insertElapse(tag string, elapse int64) {
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

//func (s *Summary) Dump(logger log.Logger)  {
//	for _, k := range s.keys {
//		logger.With("module", "main").Info("Summary", k, s.statisticMap[k])
//	}
//}

func (s *Summary) Dump(context string) {
	var res string
	for _, k := range s.keys {
		res += fmt.Sprintf("%s=%d, ", k, s.statisticMap[k])
	}
	//systemlog.Println("Elapse Summary", context, res)
	fmt.Printf("Elapse Summary: %s\n", res)
}

func (s *Summary) insert(tag string, elapse int64) {
	_, ok := s.statisticMap[tag]
	if !ok {
		return
	}
	s.statisticMap[tag] += elapse
}
