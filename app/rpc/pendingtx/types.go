package pendingtx

import "github.com/okex/exchain/x/evm/watcher"

type KafkaMsg struct {
	Topic  string               `json:"topic"`
	Source interface{}          `json:"source"`
	Data   *watcher.Transaction `json:"data"`
}
