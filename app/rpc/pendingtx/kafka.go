package pendingtx

import (
	"context"
	"encoding/json"

	rpctypes "github.com/okex/exchain/app/rpc/types"
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
		}),
	}
}

type KafkaMsg struct {
	Topic  string                `json:"topic"`
	Source interface{}           `json:"source"`
	Data   *rpctypes.Transaction `json:"data"`
}

func (kc *KafkaClient) Send(hash []byte, tx *rpctypes.Transaction) error {
	msg, err := json.Marshal(KafkaMsg{
		Topic: kc.Topic,
		Data:  tx,
	})
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
