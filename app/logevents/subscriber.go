package logevents

import (
	"fmt"
	"os"
	"time"
)

type Subscriber interface {
	Init(urls string, topic string, logdir string)
	Run()
}

func NewSubscriber() Subscriber {
	return &subscriber{
		fileMap: make(map[string]*os.File),
	}
}

type subscriber struct {
	fileMap map[string]*os.File
	kafka *logClient
	logdir string
}

func (s* subscriber) Init(urls string, topic string, logdir string)  {
	s.kafka = newLogClient(urls, HeartbeatTopic, OECLogTopic, LogConsumerGroup)
	s.logdir = logdir

	_, err := os.Stat(logdir)

	if os.IsNotExist(err) {
		err = os.Mkdir(logdir, os.ModePerm)
	}
	if err != nil {
		panic(err)
	}
}

func (s* subscriber) heartbeatRoutine() {
	ticker := time.NewTicker(HeartbeatInterval)
	pid := os.Getpid()
	id := 0
	for range ticker.C {
		key :=	fmt.Sprintf("%d:%d", pid, id)
		err := s.kafka.send(key, "heartbeat", )
		if err != nil {
			fmt.Printf("Subscriber heartbeat routine. %s, err: %s\n", key, err)
		}
		id++
		fmt.Printf("Subscriber heartbeat routine. Send: %s\n", key)
	}
}

func (s* subscriber) Run() {
	go s.heartbeatRoutine()
	for {
		key, m, err := s.kafka.recv()
		if err != nil {
			fmt.Printf("recv err: %s", err)
			continue
		}
		fmt.Printf("recv msg from %s, at topic: %v\n", key, m.Topic)
		s.onEvent(key, m.Data)
	}
}

func (s* subscriber) onEvent(from, event string)  {
	from = s.logdir + string(os.PathSeparator) + from+".log"
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
