package tx

import (
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	anytypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdktypes "github.com/okex/exchain/libs/cosmos-sdk/types"
	types "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
)

var (
	_           context.TxRequest      = (*BroadcastTxRequest)(nil)
	_           context.TxResponse     = (*BroadcastTxResponse)(nil)
	_           context.CodecSensitive = (*BroadcastTxResponse)(nil)
	nextMarshal                        = errors.New("next")
)

func (t *BroadcastTxRequest) GetModeDetail() int32 {
	return int32(t.Mode)
}

func (t *BroadcastTxRequest) GetData() []byte {
	return t.GetTxBytes()
}

func (t *BroadcastTxResponse) HandleResponse(codec *codec.CodecProxy, data interface{}) interface{} {
	resp := data.(sdktypes.TxResponse)
	logs := convOkcLogs2Proto(resp.Logs)
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
		Timestamp: resp.Timestamp,
	}

	respWrapper := AminoBroadCastTxResponse{resp}
	dataBytes, err := codec.GetCdc().MarshalJSON(respWrapper)
	if nil == err {
		bytesWrapper := &wrappers.BytesValue{Value: dataBytes}
		if any, err := anytypes.NewAnyWithValue(bytesWrapper); nil == err {
			t.TxResponse.Tx = any
		}
	}
	return t
}

func convOkcLogs2Proto(ls sdktypes.ABCIMessageLogs) types.ABCIMessageLogs {
	logs := make(types.ABCIMessageLogs, 0)
	for _, l := range ls {
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
	return logs
}

type AminoBroadCastTxResponse struct {
	TxResponse sdktypes.TxResponse `protobuf:"bytes,1,opt,name=tx_response,json=txResponse,proto3" json:"tx_response,omitempty"`
}

func (t *BroadcastTxResponse) Marshal(proxy *codec.CodecProxy) ([]byte, error) {
	if t.TxResponse == nil {
		return nil, nextMarshal
	}
	internalV := t.TxResponse.Tx.GetCachedValue()
	if nil == internalV {
		return nil, nextMarshal
	}
	resp, ok := internalV.(*wrappers.BytesValue)
	if !ok {
		return nil, nextMarshal
	}
	return resp.Value, nil
}
