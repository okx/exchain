package pendingtx

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	*kafka.Writer
}

func NewKafkaClient(addr, topic string) *KafkaClient {
	return &KafkaClient{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{addr},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (kc *KafkaClient) Send(hash, tx []byte) error {
	// Automatic retries and reconnections on errors.
	return kc.WriteMessages(context.Background(),
		kafka.Message{
			Key:   hash,
			Value: tx,
		},
	)
}
