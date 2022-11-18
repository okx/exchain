package watcher

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
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

func (aq *ACProcessorQuerier) GetAccount(key []byte) (*types.EthAccount, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgAccount); ok {
			cp := rsp.account.Copy() // the account value will be modified by caller
			return cp.(*types.EthAccount), nil
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

func (aq *ACProcessorQuerier) GetParams() (evmtypes.Params, error) {
	if v, ok := aq.p.Get(prefixParams); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return evmtypes.Params{}, nil
		}
		if rsp, ok := v.(*MsgParams); ok {
			return rsp.Params, nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds (just adaptive)
			var msgParams MsgParams
			err := json.Unmarshal(value.Value, &msgParams)
			if err != nil {
				return evmtypes.Params{}, err
			}
			return msgParams.Params, nil
		}
	}
	return evmtypes.Params{}, errACNotFound
}

// for prefixWhiteList and prefixBlackList
func (aq *ACProcessorQuerier) Has(key []byte) (bool, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return false, nil
		}
		return true, nil
	}
	return false, errACNotFound
}

func (aq *ACProcessorQuerier) GetBlackList(key []byte) ([]byte, error) {
	if v, ok := aq.p.Get(key); ok {
		if v == nil || v.GetType() == TypeDelete { // for del key
			return nil, nil
		}
		if rsp, ok := v.(*MsgContractBlockedListItem); ok {
			return []byte(rsp.GetValue()), nil
		} else if value, ok := v.(*Batch); ok { // maybe v is from the dds
			return value.Value, nil
		}
	}
	return nil, errACNotFound
}
