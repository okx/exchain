package watcher

import (
	"log"

	"github.com/golang/protobuf/proto"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/tendermint/go-amino"
)

var (
	paramsKey      = []byte("wasm-parameters")
	sendEnabledKey = []byte("send-enabled")
	supplyPrefix   = []byte("token-supply-")

	Codec *amino.Codec
)

func SetParams(para types.Params) {
	if !Enable() {
		return
	}
	b, err := proto.Marshal(&para)
	if err != nil {
		panic("wasm watchDB SetParams marshal error:" + err.Error())
	}
	if err = db.Set(paramsKey, b); err != nil {
		panic("wasm watchDB SetParams set error:" + err.Error())
	}
}

func GetParams() types.Params {
	ensureChecked()
	b, err := db.Get(paramsKey)
	if err != nil {
		panic("wasm watchDB GetParams get error:" + err.Error())
	}
	var p types.Params
	if err = proto.Unmarshal(b, &p); err != nil {
		panic("wasm watchDB GetParams unmarshal error:" + err.Error())
	}
	return p
}

type ParamsManager struct{}

func (p ParamsManager) SetSendEnabled(enable bool) {
	if !Enable() {
		return
	}
	var ok byte
	if enable {
		ok = 1
	}
	if err := db.Set(sendEnabledKey, []byte{ok}); err != nil {
		log.Println("SetSendEnabled error:", err)
	}
}

func (p ParamsManager) GetSendEnabled() bool {
	ensureChecked()
	v, err := db.Get(sendEnabledKey)
	if err != nil {
		log.Println("SetSendEnabled error:", err)
		return false
	}
	if len(v) == 0 || v[0] == 0 {
		return false
	}
	return true
}

func (p ParamsManager) SetSupply(supply interface{}) {
	if !Enable() {
		return
	}
	tokensSupply, ok := supply.(sdk.Coins)
	if !ok {
		return
	}

	for i := range tokensSupply {
		err := db.Set(getTokenSupplyKey(tokensSupply[i].Denom), Codec.MustMarshalBinaryLengthPrefixed(tokensSupply[i].Amount))
		if err != nil {
			log.Printf("wasm watchDB SetSupply, token: %s, error: %s\n", tokensSupply[i].Denom, err)
		}
	}
}

func (p ParamsManager) GetSupply() interface{} {
	ensureChecked()
	start := supplyPrefix
	end := cpIncr(supplyPrefix)
	iter, err := db.Iterator(start, end)
	if err != nil {
		log.Println("GetSupply Iterator error:", err)
		return sdk.NewCoins()
	}
	defer iter.Close()
	var coins sdk.Coins
	for ; iter.Valid(); iter.Next() {
		var amount sdk.Dec
		Codec.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &amount)
		coins = append(coins, sdk.NewCoin(string(iter.Key()[len(supplyPrefix):]), amount))
	}
	return coins
}

func getTokenSupplyKey(demon string) []byte {
	return append(supplyPrefix, demon...)
}

func cpIncr(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {
			// Overflow
			return nil
		}
	}
	return nil
}
