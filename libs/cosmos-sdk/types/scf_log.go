package types

import "fmt"

var (
	DebugLogByScf = NewScfLog()
)

func NewScfLog() *ScfLog {
	return &ScfLog{
		rwSet:      make([]string, 0),
		commitInfo: make([]string, 0),
	}
}

type ScfLog struct {
	rwSet      []string
	commitInfo []string
}

func (s *ScfLog) Clean() {
	s.commitInfo = make([]string, 0)
	s.rwSet = make([]string, 0)
}

func (s *ScfLog) AddCommitInfo(data string) {
	s.commitInfo = append(s.commitInfo, data)
}

func (s *ScfLog) AddRWSet(data []string) {
	s.rwSet = append(s.rwSet, data...)
}

func (s *ScfLog) PrintDebugInfo() {
	fmt.Println("begin print commit info")
	for _, v := range s.commitInfo {
		fmt.Println(v)
	}

	fmt.Println("detail rwset")
	for _, v := range s.rwSet {
		fmt.Println(v)
	}

}
