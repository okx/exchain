package tx

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/types"
	types "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

var (
	_ context.TxRequest  = (*BroadcastTxRequest)(nil)
	_ context.TxResponse = (*BroadcastTxResponse)(nil)
)

func (t *BroadcastTxRequest) GetModeDetail() int32 {
	return int32(t.Mode)
}

func (t *BroadcastTxRequest) GetData() []byte {
	return t.GetTxBytes()
}

func (t *BroadcastTxResponse) HandleResponse(data interface{}) interface{} {
	resp := data.(types2.TxResponse)
	fmt.Println(resp)
	logs := make(types.ABCIMessageLogs, 0)
	for _, l := range resp.Logs {
		es := make(types.StringEvents, 0)
		for _, e := range l.Events {
			attrs := make([]types.Attribute, 0)
			for _, a := range e.Attributes {
				attrs = append(attrs, types.Attribute{
					Key:   a.Key,
					Value: a.Value,
				})
			}
			es = append(es, types.StringEvent{
				Type:       e.Type,
				Attributes: attrs,
			})
		}
		logs = append(logs, types.ABCIMessageLog{
			MsgIndex: uint32(l.MsgIndex),
			Log:      l.Log,
			Events:   es,
		})
	}
	t.TxResponse = &types.TxResponse{
		Height:    resp.Height,
		TxHash:    resp.TxHash,
		Codespace: resp.Codespace,
		Code:      resp.Code,
		Data:      resp.Data,
		RawLog:    resp.RawLog,
		Logs:      logs,
		Info:      resp.Info,
		GasWanted: resp.GasWanted,
		GasUsed:   resp.GasUsed,
		//Tx:        resp.Tx,
		Timestamp: resp.Timestamp,
	}
	return t
}
