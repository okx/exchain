package watcher

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/app/types"
	prototypes "github.com/okex/exchain/x/evm/watcher/proto"
	"strconv"
)

var errACNotFound = errors.New("ac processor: not found")

type ACProcessorQuerier struct {
	p *ACProcessor
}

func newACProcessorQuerier(p *ACProcessor) *ACProcessorQuerier {
	if p != nil {
		return &ACProcessorQuerier{p: p}
	}
	// local
	return &ACProcessorQuerier{p: &ACProcessor{
		commitList:  newCommitCache(), // for support to querier
		curMsgCache: newMessageCache(),
	}}
}

func (aq *ACProcessorQuerier) GetTransactionReceipt(key []byte) (*TransactionReceipt, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if value, ok := v.(*MsgTransactionReceipt); ok {
			value.ObjectParse()
			return value.TransactionReceipt, nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			var protoReceipt prototypes.TransactionReceipt
			e := proto.Unmarshal(value.Value, &protoReceipt)
			if e != nil {
				return nil, e
			}
			receipt := protoToReceipt(&protoReceipt)
			return receipt, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetTransactionResponse(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if value, ok := v.(*MsgStdTransactionResponse); ok {
			return []byte(value.txResponse), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetBlockByHash(key []byte) (*Block, error) {
	var b []byte
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgBlock); ok {
			b = []byte(rsp.block)
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			b = value.Value
		} else {
			return nil, errACNotFound
		}
		var block Block
		err := json.Unmarshal(b, &block)
		if err != nil {
			return nil, err
		}
		return &block, nil
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetTransactionByHash(key []byte) (*Transaction, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgEthTx); ok {
			err := rsp.ObjectParse()
			if err != nil {
				return nil, err
			}
			return rsp.Transaction, nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			var protoTx prototypes.Transaction
			e := proto.Unmarshal(value.Value, &protoTx)
			if e != nil {
				return nil, e
			}
			return protoToTransaction(&protoTx), nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetLatestBlockNumber(key []byte) (uint64, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return 0, nil
		}
		var height []byte
		if rsp, ok := v.(*MsgLatestHeight); ok {
			height = []byte(rsp.GetValue())
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			height = value.Value
		}
		h, e := strconv.Atoi(string(height))
		return uint64(h), e
	}
	return 0, errACNotFound
}

func (aq *ACProcessorQuerier) GetBlockHashByNumber(key []byte) (common.Hash, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return common.Hash{}, nil
		}
		if rsp, ok := v.(*MsgBlockInfo); ok {
			return common.HexToHash(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return common.HexToHash(string(value.Value)), nil
		}
	}
	return common.Hash{}, errACNotFound
}

func (aq *ACProcessorQuerier) GetBlockHash(key []byte) (common.Hash, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return common.Hash{}, nil
		}
		if rsp, ok := v.(*MsgBlockInfo); ok {
			return common.HexToHash(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return common.HexToHash(string(value.Value)), nil
		}
	}
	return common.Hash{}, errACNotFound
}

func (aq *ACProcessorQuerier) GetCode(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgCode); ok {
			return []byte(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetCodeByHash(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgCodeByHash); ok {
			return []byte(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetStdTxHash(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgBlockStdTxHash); ok {
			return []byte(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetAccount(key []byte) (*types.EthAccount, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgAccount); ok {
			return rsp.account, nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			acc, err := DecodeAccount(value.Value)
			if err != nil {
				return nil, err
			}
			return acc, nil
		}
	}
	return nil, errACNotFound
}

func (aq *ACProcessorQuerier) GetState(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgState); ok {
			return rsp.value, nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}
