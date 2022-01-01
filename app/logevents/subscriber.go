package logevents

import "os"

type Subscriber interface {
	Init(urls string, topic string)
	Run()
}

func NewSubscriber() Subscriber {
	return &subscriber{
		fileMap: make(map[string]*os.File),
	}
}

type subscriber struct {
	fileMap map[string]*os.File
}

func (s* subscriber) Init(urls string, topic string)  {

}

func (s* subscriber) Run() {

}

func (s* subscriber) onEvent(from, event string)  {
	f, err := s.getOsFile(from)
	if err != nil {
		return
	}

	_, err = f.WriteString(event)
	if err != nil {
		return
	}
}


func (s* subscriber) getOsFile(fileName string) (file *os.File, err error) {
	var ok bool
	file, ok = s.fileMap[fileName]

	if ok {
		return
	}

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		s.fileMap[fileName] = file
	}
	return
}
