package evm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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
	RegisterHandle()
}

func RegisterHandle() {
	baseapp.RegisterEvmResultConverter(EncodeResultData)
	baseapp.RegisterEvmConvertJudge(EvmConvertJudge)
}

func EvmConvertJudge(msg sdk.Msg) (*baseapp.CMTxParam, []byte, bool) {
	if msg.Route() != types.ModuleName {
		return nil, nil, false
	}
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if !ok || evmTx.Data.Recipient == nil { // deploy contract no need convert to cosmos msg
		return nil, nil, false
	}
	cmtp, err := ContractStringParamParse(evmTx.Data.Payload)
	if err != nil {
		return nil, nil, false
	}
	return cmtp, evmTx.Data.Recipient[:], true
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
