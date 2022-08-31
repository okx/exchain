package pendingtx

import (
	"context"

	"github.com/okex/exchain/x/evm/watcher"
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Topic string
	*kafka.Writer
}

func NewKafkaClient(addrs []string, topic string) *KafkaClient {
	return &KafkaClient{
		Topic: topic,
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  addrs,
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		}),
	}
}

func (kc *KafkaClient) Send(hash []byte, tx *watcher.Transaction) error {
	kafkaMsg := KafkaMsg{
		Topic: kc.Topic,
		Data:  tx,
	}

	msg, err := kafkaMsg.MarshalJSON()
	if err != nil {
		return err
	}

	// Automatic retries and reconnections on errors.
	return kc.WriteMessages(context.Background(),
		kafka.Message{
			Key:   hash,
			Value: msg,
		},
	)
}
