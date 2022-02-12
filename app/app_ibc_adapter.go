package app

import (
	//"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	types2 "github.com/okex/exchain/libs/tendermint/abci/types"
	//sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	//"github.com/okex/exchain/x/evm"
)

func MakeIBC() types.InterfaceRegistry {
	interfaceReg := types.NewInterfaceRegistry()
	ModuleBasics.RegisterInterfaces(interfaceReg)
	return interfaceReg
}

//func IBCMergeDecoder(cdc *codec.Codec, m codec.Marshaler) sdk.TxDecoder {
//	return func(txBytes []byte, height ...int64) (sdk.Tx, error) {
//		ret,err:=evm.TxDecoder(cdc)(txBytes,height...)
//		if nil==err && ret!=nil{
//			return ret,nil
//		}
//		var msg sdk.MsgAdapter
//		m.UnmarshalBinaryBare()
//	}
//}

func (a *OKExChainApp) Upgrade(ctx sdk.Context, req *types2.UpgradeReq) (*types2.UpgradeResp, error) {
	resp, err := a.mm.Upgrade(req)
	if nil != err {
		panic(err)
	}
	return resp, nil
}
