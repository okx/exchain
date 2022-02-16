package logevents

import (
	"fmt"
	"github.com/okex/exchain/libs/system"
	"os"
	"time"
)

type Subscriber interface {
	Init(urls string, logdir string)
	Run()
}

func NewSubscriber() Subscriber {
	return &subscriber{
		fileMap: make(map[string]*os.File),
	}
}

type subscriber struct {
	fileMap map[string]*os.File
	kafka   *logClient
	logdir  string
}

func (s *subscriber) Init(urls string, logdir string) {
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

func (s *subscriber) heartbeatRoutine() {
	ticker := time.NewTicker(HeartbeatInterval)
	pid := system.Getpid()
	id := 0
	for range ticker.C {
		key := fmt.Sprintf("%d:%d", pid, id)
		msg := &KafkaMsg{
			Data: "heartbeat",
		}
		err := s.kafka.send(key, msg)
		if err != nil {
			fmt.Printf("Subscriber heartbeat routine. %s, err: %s\n", key, err)
			continue
		}
		id++
		fmt.Printf("Subscriber heartbeat routine. Send: %s\n", key)
	}
}

func (s *subscriber) Run() {
	go s.heartbeatRoutine()
	for {
		key, m, err := s.kafka.recv()
		if err != nil {
			fmt.Printf("recv err: %s", err)
			continue
		}
		fmt.Printf("recv msg from %s, at topic: %v\n", key, m.Topic)
		err = s.onEvent(key, m.Data)
		if err != nil {
			fmt.Printf("onEvent err: %s", err)
		}
	}
}

func (s *subscriber) onEvent(from, event string) (err error) {
	from = s.logdir + string(os.PathSeparator) + from + ".log"

	var f *os.File
	f, err = s.getOsFile(from)
	if err != nil {
		return
	}

	_, err = f.WriteString(event)
	return
}

func (s *subscriber) getOsFile(fileName string) (file *os.File, err error) {
	var ok bool
	file, ok = s.fileMap[fileName]

	if ok {
		return
	}

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err == nil {
		s.fileMap[fileName] = file
	}
	return
}
