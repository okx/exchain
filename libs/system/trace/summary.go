package trace

var sum *Summary
type Summary struct {
	statisticMap map[string]int64
}

func init() {
	sum = &Summary{
		statisticMap: make(map[string]int64),
	}
	//analyzer.SetInsertFunc(InsertElapse)
}

func InsertElapse(tag string, elapse int64)  {
	sum.Insert(tag, elapse)
}

func GetTraceSummary() *Summary {
	return sum
}

func (s *Summary) Init()  {

}

func (s *Summary) Insert(tag string, elapse int64)  {
	_, ok := s.statisticMap[tag]
	if !ok {
		return
	}
	s.statisticMap[tag] += elapse
}