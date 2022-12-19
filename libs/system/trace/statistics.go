package trace

import "time"

type BaseStatistics interface {
	Init(tags ...string)
	Accumulate(tag string, lastPinTime time.Time)
	GetTags() []string
	GetValue(tag string) int64
}

type Summary struct {
	statisticMap map[string]int64
	keys         []string
}

func NewSummary() *Summary {
	return &Summary{
		statisticMap: make(map[string]int64),
	}
}

func (s *Summary) Init(tags ...string) {
	for _, k := range tags {
		s.statisticMap[k] = 0
	}
	s.keys = tags
}

func (s *Summary) Accumulate(tag string, lastPinTime time.Time) {
	s.statisticMap[tag] += time.Since(lastPinTime).Nanoseconds()
}

func (s *Summary) GetTags() []string {
	return s.keys
}

func (s *Summary) GetValue(tag string) int64 {
	return s.statisticMap[tag]
}

type StatisticsCell interface {
	StartTiming()
	EndTiming(tag string)
}
