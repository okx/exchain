package types

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/tendermint/go-amino"
)

type CM40EventDataNewBlock struct {
	Block *CM40Block `json:"block"`

	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
}

func RegisterCM40EventDatas(cdc *amino.Codec) {
	registerCommonEventDatas(cdc)
	cdc.RegisterConcrete(CM40EventDataNewBlock{}, "tendermint/event/NewBlock", nil)
}
