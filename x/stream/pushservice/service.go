package pushservice

import (
	"encoding/json"
	"fmt"

	"github.com/okex/exchain/x/stream/pushservice/conn"

	"github.com/okex/exchain/x/backend"

	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/x/stream/pushservice/types"
	"github.com/okex/exchain/x/token"
)

var _ types.Writer = (*PushService)(nil)

// PushService is to push data to redis_push_service_channel
type PushService struct {
	client *conn.Client
	log    log.Logger
}

func NewPushService(redisURL, redisPassword string, db int, log log.Logger) (srv *PushService, err error) {
	c, err := conn.NewClient(redisURL, redisPassword, db, log)
	if err != nil {
		log.Error("connect pushservice",
			"err", err.Error(),
		)
		return nil, err
	}
	log = log.With("module", "pushservice")

	return &PushService{client: c, log: log}, nil
}

// WriteSync push data to redis-push-channel per block
func (p PushService) WriteSync(b *types.RedisBlock) (map[string]int, error) {
	result := make(map[string]int)

	// orders
	for _, val := range b.OrdersMap {
		result["orders"] += len(val)
	}
	for k, v := range b.OrdersMap {
		if err := p.setOrders(k, v); err != nil {
			return result, fmt.Errorf("setOrders failed, %s", err.Error())
		}
	}

	// accounts
	result["accs"] = len(b.AccountsMap)
	for k, v := range b.AccountsMap {
		if err := p.setAccount(k, v); err != nil {
			return result, fmt.Errorf("setAccount failed, %s", err.Error())
		}
	}

	// match results
	result["matches"] += len(b.MatchesMap)
	for k, v := range b.MatchesMap {
		if err := p.setMatches(k, v); err != nil {
			return result, fmt.Errorf("setMatches failed, %s", err.Error())
		}
	}

	// product depth
	for _, val := range b.DepthBooksMap {
		result["depth"] += len(val.Asks) + len(val.Bids)
	}
	for k, v := range b.DepthBooksMap {
		if err := p.setDepthSnapshot(k, v); err != nil {
			return result, fmt.Errorf("setDepthSnapshot failed, %s", err.Error())
		}
	}

	// coins, products
	result["instruments"] = len(b.Instruments)
	if len(b.Instruments) != 0 {
		instrs := make([]string, 0, len(b.Instruments))
		for k := range b.Instruments {
			instrs = append(instrs, k)
		}
		if err := p.setInstruments(instrs); err != nil {
			return result, fmt.Errorf("setInstruments failed, %s", err.Error())
		}
	}

	b.Clear()
	return result, nil
}

// Close connection to remote redis server
func (p PushService) Close() error {
	return p.client.Close()
}

// get redis client
func (p *PushService) GetConnCli() *conn.Client {
	return p.client
}

// setAccount, push account to private channel
func (p PushService) setAccount(address string, info token.CoinInfo) error {
	value, err := json.Marshal(info)
	if err != nil {
		return err
	}
	key := address
	p.log.Debug("setAccount", "key", key, "value", string(value))
	return p.client.PrivatePub(key, string(value))
}

// setOrders push orders to private channel
func (p PushService) setOrders(address string, orders []backend.Order) error {
	value, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	key := address
	p.log.Debug("setOrders", "key", key, "value", string(value))
	return p.client.PrivatePub(key, string(value)[1:len(string(value))-1])
}

// setMatches push matches to public channel
func (p PushService) setMatches(product string, matches backend.MatchResult) error {
	value, err := json.Marshal(matches)
	if err != nil {
		return err
	}
	key1 := product
	key2 := product
	p.log.Debug("setMatches_pub", "key", key1, "value", string(value))
	p.log.Debug("setMatches_set", "key", key2, "value", string(value))
	if err := p.client.Set(key2, string(value)); err != nil {
		return err
	}
	return p.client.PublicPub(key1, string(value))
}

// setDepthSnapshot push depth to public channel
func (p PushService) setDepthSnapshot(product string, depth types.BookRes) error {
	value, err := json.Marshal(depth)
	if err != nil {
		return err
	}
	key1 := product
	key2 := product
	p.log.Debug("setDepthSnapshot_pub", "key", key1, "value", string(value))
	p.log.Debug("setDepthSnapshot_set", "key", key2, "value", string(value))
	if err := p.client.Set(key2, string(value)); err != nil {
		return err
	}
	return p.client.DepthPub(key1, string(value))
}

// setInstruments push instruments to
func (p PushService) setInstruments(instruments []string) error {
	value, err := json.Marshal(instruments)
	if err != nil {
		return err
	}
	key := "instruments"
	p.log.Debug("setInstruments", "key", key, "value", string(value))
	return p.client.Set(key, string(value))
}
