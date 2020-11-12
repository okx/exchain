package kafkaclient

import (
	"encoding/binary"
	"github.com/segmentio/kafka-go"
)

type OKDExBalancer struct {
	maxPartition int64
}

func newOKDExBalancer(maxPartition int64) OKDExBalancer {
	return OKDExBalancer{
		maxPartition: maxPartition,
	}
}

func (ob OKDExBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	return int(int64(binary.BigEndian.Uint64(msg.Key)) % ob.maxPartition)
}