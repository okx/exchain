package pendingtx

import (
	"context"
	"encoding/json"

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
			Balancer: &kafka.Hash{},
			Async:    true,
		}),
	}
}

func (kc *KafkaClient) SendPending(hash []byte, tx *PendingTx) error {
	kafkaMsg := PendingMsg{
		Topic: kc.Topic,
		Data:  tx,
	}

	msg, err := json.Marshal(kafkaMsg)
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

func (kc *KafkaClient) SendRmPending(hash []byte, tx *RmPendingTx) error {
	kafkaMsg := RmPendingMsg{
		Topic: kc.Topic,
		Data:  tx,
	}

	msg, err := json.Marshal(kafkaMsg)
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
