package logevents

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"strings"
	"sync"
	"time"
)

const (
	OECLogTopic      = "oeclog"
	LogConsumerGroup = "oeclog-consumer-group"

	HeartbeatTopic = "oeclog-subscriber-heartbeat"

	HeartbeatInterval = 5 * time.Second
	ExpiredInterval   = 6 * HeartbeatInterval
)

type logClient struct {
	wt      string
	rt      string
	groupID string
	*kafka.Writer
	*kafka.Reader
}

func newLogClient(kafkaAddrs string, wt, rt string, groupID string) *logClient {
	addrs := strings.Split(kafkaAddrs, ",")
	return &logClient{
		wt:      wt,
		rt:      rt,
		groupID: groupID,
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  addrs,
			Topic:    wt,
			Balancer: &kafka.LeastBytes{},
		}),
		Reader: getKafkaReader(kafkaAddrs, rt, groupID),
	}
}

type KafkaMsg struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var KafkaMsgPool = sync.Pool{
	New: func() interface{} {
		return &KafkaMsg{}
	},
}

func (kc *logClient) recv() (string, *KafkaMsg, error) {
	const empty = ""
	rawMsg, err := kc.ReadMessage(context.Background())
	if err != nil {
		return empty, nil, err
	}

	var msg KafkaMsg
	err = json.Unmarshal(rawMsg.Value, &msg)
	if err != nil {
		return empty, nil, err
	}

	return string(rawMsg.Key), &msg, err
}

func (kc *logClient) send(key string, rawMsg *KafkaMsg) error {
	rawMsg.Topic = kc.wt

	msg, err := json.Marshal(*rawMsg)
	if err != nil {
		return err
	}

	// Automatic retries and reconnections on errors.
	return kc.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: msg,
		},
	)
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
		//MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
