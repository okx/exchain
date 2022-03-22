package app

import (
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

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

func (app *OKExChainApp) RegisterTxService(clientCtx clientCtx.CLIContext) {

}
