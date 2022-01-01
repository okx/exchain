package logevents

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/spf13/viper"
	"net"
	"strings"
)


const (
	FlagLogServerUrl string = "log-server"
)


type provider struct {
	eventChan chan string
	identity string
	logServerUrl string
	logger log.Logger
}

func NewProvider(logger log.Logger) log.Subscriber {
	url := viper.GetString(FlagLogServerUrl)
	if len(url) == 0 {
		return nil
	}

	p := &provider{
		eventChan: make(chan string, 1000),
		logServerUrl: url,
		logger: logger.With("module", "provider"),
	}
	p.init()
	return p
}

func (p* provider) init()  {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		var comma string
		for _, value := range addrs {
			if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					p.identity += fmt.Sprintf("%s%s", comma, ipnet.IP.String())
					break
				}
			}
		}
	}
	if len(p.identity) == 0 {
		panic("")
	}

	p.logger.Info("provider init",
		"url", p.logServerUrl,
		"id", p.identity)

	go p.eventRoutine()
}

func (p* provider) AddEvent(event string)  {
	if strings.Index(event, "module=provider") != -1 {
		return
	}
	p.eventChan <- event
}

func (p* provider) eventRoutine()  {
	for event := range p.eventChan {
		p.eventHandler(event)
	}
}

func (p* provider) eventHandler(event string)  {
	p.logger.Info("new event", "event", event)
	//write("sub.txt", event)
}

//func write(file, param string) {
//	//f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
//	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
//	if err != nil {
//		return
//	}
//	defer f.Close()
//
//	_, err = f.WriteString(param)
//	if err != nil {
//		return
//	}
//}
