package evm2cm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
)

var (
	ErrInputDataSize = errors.New("the input data size is error")
)

func init() {
	RegisterEvmHandle()
}

func RegisterEvmHandle() {
	baseapp.RegisterEvmResultConverter(EncodeResultData)
	baseapp.RegisterEvmMsgParser(EvmMsgParser)
	baseapp.RegisterEvmConvertJudge(EvmConvertJudge)
}

func EvmConvertJudge(filter map[string]struct{}, msg sdk.Msg) bool {
	if msg.Route() != types.ModuleName {
		return false
	}
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if !ok || evmTx.Data.Recipient == nil { // deploy contract no need convert to cosmos msg
		return false
	}
	to := evmTx.Data.Recipient.String()
	if filter != nil {
		if _, ok := filter[to]; ok {
			return true
		}
	}
	return false
}

func EvmMsgParser(msg sdk.Msg) (*baseapp.CMTxParam, error) {
	if evmTx, ok := msg.(*types.MsgEthereumTx); ok {
		return ContractStringParamParse(evmTx.Data.Payload)
	}
	return nil, fmt.Errorf("msg is not a MsgEthereumTx")
}

func ContractStringParamParse(input []byte) (*baseapp.CMTxParam, error) {
	const methodSite = 4
	const fixedSite = 32
	const padSite = methodSite + fixedSite                 // 36
	const dataLenSite = methodSite + fixedSite + fixedSite // 68
	if len(input) < dataLenSite {
		return nil, ErrInputDataSize
	}

	size := new(big.Int).SetBytes(input[padSite:dataLenSite]) // 存放数据长度
	if len(input) < int(size.Int64())+dataLenSite {
		return nil, ErrInputDataSize
	}
	data := input[dataLenSite : size.Int64()+dataLenSite] // 实际数据

	value, err := hex.DecodeString(string(data)) // this is json fmt
	if err != nil {
		return nil, err
	}
	cmtx := &baseapp.CMTxParam{}
	err = json.Unmarshal(value, cmtx)
	if err != nil {
		return nil, err
	}
	return cmtx, nil
}

func EncodeResultData(data []byte) ([]byte, error) {
	ethHash := common.BytesToHash(data)
	return types.EncodeResultData(&types.ResultData{TxHash: ethHash})
}
