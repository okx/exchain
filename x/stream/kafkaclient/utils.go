package kafkaclient

import "encoding/binary"

func getKafkaMsgKey(marketID int64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(marketID))
	return key
}
