package tx

import (
	"encoding/hex"
	"testing"

	okexchaincodec "github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	ibctransfer "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	"github.com/stretchr/testify/require"
)

func TestIbcTxDecoderSuccess(t *testing.T) {
	ModuleBasics := module.NewBasicManager(
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
	)
	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)
	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	cpcdc := &codec.CompoundCodec{
		cdc,
		codecProxy,
	}

	txBytes, err := hex.DecodeString("0a8c030a89030a232f6962632e636f72652e636c69656e742e76312e4d7367437265617465436c69656e7412e1020aaa010a2b2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436c69656e745374617465127b0a056962632d311204080110031a040880ac4d22040880df6e2a0308d80432003a040801102142190a090801180120012a0100120c0a02000110211804200c300142190a090801180120012a0100120c0a02000110201801200130014a07757067726164654a1075706772616465644942435374617465500158011286010a2e2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436f6e73656e737573537461746512540a0c08a2eb9992061080a4b1b60312220a206c4821c63adcb62c57effb912d6d568f8971456acf6630c5b3677c8d3cd999321a203ee89cb19deef923d4ffea8d44980336746dd0125d7372f9a78562a8aaf0fb821a296578316865756b7866766d74617968306a7a746468706e73756468677a666d3364783772357261753412370a290a210a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912040a020801120a0a080a036f6b741201301a00")
	require.NoError(t, err)

	marshaler := cpcdc.GetProtocMarshal()
	decode := IbcTxDecoder(marshaler)
	_, err = decode(txBytes)

	require.NoError(t, err)
}

func TestIbcTxDecoderFailure(t *testing.T) {
	ModuleBasics := module.NewBasicManager(
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
	)
	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)
	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	cpcdc := &codec.CompoundCodec{
		cdc,
		codecProxy,
	}

	txBytes, err := hex.DecodeString("7261753412370a290a210a1f2f636f736d6f732e637")

	marshaler := cpcdc.GetProtocMarshal()
	decode := IbcTxDecoder(marshaler)
	_, err = decode(txBytes)
	require.EqualError(t, err, "tx parse error: invalid length; unexpected EOF")
}
