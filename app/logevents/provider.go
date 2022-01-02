package logevents

import (
	"bytes"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/system"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type provider struct {
	eventChan chan *KafkaMsg
	identity string
	logServerUrl string
	logger log.Logger
	kafka *logClient
	subscriberAlive bool

	mutex sync.Mutex
	lastHeartbeat time.Time
}

func NewProvider(logger log.Logger) log.Subscriber {
	url := viper.GetString(server.FlagLogServerUrl)
	if len(url) == 0 {
		logger.Info("Publishing logs is disabled")
		return nil
	}

	p := &provider{
		eventChan: make(chan *KafkaMsg, 1000),
		logServerUrl: url,
		logger: logger.With("module", "provider"),
	}
	p.init()
	return p
}

func (p* provider) init()  {

	var err error
	p.identity, err = system.GetIpAddr(viper.GetBool(types.FlagAppendPid))

	if len(p.identity) == 0 {
		panic("Invalid identity")
	}

	if err != nil{
		p.logger.Error("Failed to set identity", "err", err)
		return
	}

	role := viper.GetString("consensus-role")
	if len(role) > 0 {
		p.identity = role
	}

	p.kafka = newLogClient(p.logServerUrl, OECLogTopic, HeartbeatTopic, p.identity)

	p.logger.Info("Provider init", "url", p.logServerUrl, "id", p.identity)

	go p.eventRoutine()
	go p.expiredRoutine()
	go p.heartbeatRoutine()
}

func (p* provider) AddEvent(buf *bytes.Buffer)  {
	if !p.subscriberAlive {
		return
	}

	msg := KafkaMsgPool.Get().(*KafkaMsg)
	msg.Data = buf.String()
	p.eventChan <- msg
}

func (p* provider) eventRoutine()  {
	for event := range p.eventChan {
		p.eventHandler(event)
	}
}

func (p* provider) heartbeatInterval() time.Duration {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return time.Now().Sub(p.lastHeartbeat)
}

func (p* provider) keepAlive() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.lastHeartbeat = time.Now()
	p.subscriberAlive = true
}

func (p* provider) expiredRoutine() {
	ticker := time.NewTicker(ExpiredInterval)
	for range ticker.C {
		interval := p.heartbeatInterval()
		if interval > ExpiredInterval {
			p.subscriberAlive = false
			p.logger.Info("Subscriber expired", "not-seen-for", interval, )
		}
	}
}

func (p* provider) heartbeatRoutine()  {
	for {
		key, m, err := p.kafka.recv()
		if err != nil {
			p.logger.Error("Provider heartbeat routine", "err", err)
			continue
		}
		p.logger.Info("Provider heartbeat routine",
			"from", key,
			"value", m.Data,
			//"topic", m.Topic,
			"err", err,
			)
		p.keepAlive()
	}
}

func (p* provider) eventHandler(msg *KafkaMsg)  {
	// DO NOT use p.logger to log anything in this method!!!
	defer KafkaMsgPool.Put(msg)
	p.kafka.send(p.identity, msg)
}
