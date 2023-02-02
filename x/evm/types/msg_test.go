package types

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	ibcfee "github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	okexchaincodec "github.com/okex/exchain/app/codec"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	ibctxdecode "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	ibctransfer "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
)

func newSdkAddress() sdk.AccAddress {
	tmpKey := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(tmpKey.Address().Bytes())
}

func TestMsgEthereumTx(t *testing.T) {
	addr := GenerateEthAddress()

	msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))
	require.NotNil(t, msg)
	require.Equal(t, *msg.Data.Recipient, addr)
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgEthereumTx)
	require.NotNil(t, msg.To())
	require.Equal(t, msg.GetMsgs(), []sdk.Msg{msg})
	require.Panics(t, func() { msg.GetSigners() })
	require.Panics(t, func() { msg.GetSignBytes() })
	require.Equal(t, msg.GetNonce(), uint64(0))

	msg = NewMsgEthereumTxContract(0, nil, 100000, nil, []byte("test"))
	require.NotNil(t, msg)
	require.Nil(t, msg.Data.Recipient)
	require.Nil(t, msg.To())

}

func TestTxFnSignatureInfo(t *testing.T) {
	type expected struct {
		sig []byte
		i   int
	}
	testCases := []struct {
		msg      string
		fn       func() *MsgEthereumTx
		expected expected
	}{
		{
			"receipt not nil should equal",
			func() *MsgEthereumTx {
				addr := ethcmn.BytesToAddress([]byte("test_address"))
				msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))
				return msg
			},
			expected{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74},
				0,
			},
		},
		{
			"payload below 4 bytes should DefaultSendCoinFnSignature",
			func() *MsgEthereumTx {
				addr := ethcmn.BytesToAddress([]byte("test_address"))
				msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("t"))
				return msg
			},
			expected{
				DefaultSendCoinFnSignature,
				0,
			},
		},
		{
			"receipt nil should be DefaultDeployContractFnSignature",
			func() *MsgEthereumTx {
				msg := NewMsgEthereumTx(0, nil, nil, 100000, nil, []byte("t"))
				return msg
			},
			expected{
				DefaultDeployContractFnSignature,
				1,
			},
		},
	}
	for _, tc := range testCases {
		msg := tc.fn()
		r, i := msg.GetTxFnSignatureInfo()
		require.Equal(t, tc.expected.i, i)
		require.Equal(t, tc.expected.sig, r)
	}
}

func TestMsgEthereumTxValidation(t *testing.T) {
	testCases := []struct {
		msg        string
		amount     *big.Int
		gasPrice   *big.Int
		expectPass bool
	}{
		{msg: "pass", amount: big.NewInt(100), gasPrice: big.NewInt(100000), expectPass: true},
		{msg: "pass amount is zero", amount: big.NewInt(0), gasPrice: big.NewInt(100000), expectPass: true},
		{msg: "invalid amount", amount: big.NewInt(-1), gasPrice: big.NewInt(100000), expectPass: false},
		{msg: "invalid gas price", amount: big.NewInt(100), gasPrice: big.NewInt(-1), expectPass: false},
		{msg: "invalid gas price", amount: big.NewInt(100), gasPrice: big.NewInt(0), expectPass: false},
	}

	for i, tc := range testCases {
		msg := NewMsgEthereumTx(0, nil, tc.amount, 0, tc.gasPrice, nil)

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "valid test %d failed: %s", i, tc.msg)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "invalid test %d passed: %s", i, tc.msg)
		}
	}
}

func TestMsgEthereumTxRLPSignBytes(t *testing.T) {
	addr := ethcmn.BytesToAddress([]byte("test_address"))
	chainID := big.NewInt(3)

	msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))
	hash := msg.RLPSignBytes(chainID)
	require.Equal(t, "5BD30E35AD27449390B14C91E6BCFDCAADF8FE44EF33680E3BC200FC0DC083C7", fmt.Sprintf("%X", hash))
}

func TestMsgEthereumTxRLPEncode(t *testing.T) {
	addr := ethcmn.BytesToAddress([]byte("test_address"))
	msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))

	raw, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	require.Equal(t, ethcmn.FromHex("E48080830186A0940000000000000000746573745F61646472657373808474657374808080"), raw)
}

func TestMsgEthereumTxRLPDecode(t *testing.T) {
	var msg MsgEthereumTx

	raw := ethcmn.FromHex("E48080830186A0940000000000000000746573745F61646472657373808474657374808080")
	addr := ethcmn.BytesToAddress([]byte("test_address"))
	expectedMsg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))

	err := rlp.Decode(bytes.NewReader(raw), &msg)
	require.NoError(t, err)
	require.Equal(t, expectedMsg.Data, msg.Data)

	// value size exceeds available input length of stream
	mockStream := rlp.NewStream(bytes.NewReader(raw), 1)
	require.Error(t, msg.DecodeRLP(mockStream))
}

func TestMsgEthereumTxSig(t *testing.T) {
	chainID, zeroChainID := big.NewInt(3), big.NewInt(0)

	priv1, _ := ethsecp256k1.GenerateKey()
	priv2, _ := ethsecp256k1.GenerateKey()
	addr1 := ethcmn.BytesToAddress(priv1.PubKey().Address().Bytes())
	trimed := strings.TrimPrefix(addr1.Hex(), "0x")

	fmt.Printf("%s\n", trimed)
	addrSDKAddr1, err := sdk.AccAddressFromHex(trimed)
	require.NoError(t, err)
	addr2 := ethcmn.BytesToAddress(priv2.PubKey().Address().Bytes())

	// require valid signature passes validation
	msg := NewMsgEthereumTx(0, &addr1, nil, 100000, nil, []byte("test"))
	err = msg.Sign(chainID, priv1.ToECDSA())
	require.Nil(t, err)

	err = msg.VerifySig(chainID, 0)
	require.NoError(t, err)
	require.Equal(t, addr1, msg.EthereumAddress())
	require.NotEqual(t, addr2, msg.EthereumAddress())

	signers := msg.GetSigners()
	require.Equal(t, 1, len(signers))
	require.True(t, addrSDKAddr1.Equals(signers[0]))

	// zero chainID
	err = msg.Sign(zeroChainID, priv1.ToECDSA())
	require.Nil(t, err)
	err = msg.VerifySig(zeroChainID, 0)
	require.Nil(t, err)

	// require invalid chain ID fail validation
	msg = NewMsgEthereumTx(0, &addr1, nil, 100000, nil, []byte("test"))
	err = msg.Sign(chainID, priv1.ToECDSA())
	require.Nil(t, err)
}

func TestMsgEthereumTx_ChainID(t *testing.T) {
	chainID := big.NewInt(3)
	priv, _ := ethsecp256k1.GenerateKey()
	addr := ethcmn.BytesToAddress(priv.PubKey().Address().Bytes())
	msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))
	err := msg.Sign(chainID, priv.ToECDSA())
	require.Nil(t, err)

	require.True(t, chainID.Cmp(msg.ChainID()) == 0)

	msg.Data.V = big.NewInt(27)
	require.NotNil(t, msg.ChainID())

	msg.Data.V = math.MaxBig256
	expectedChainID := new(big.Int).Div(new(big.Int).Sub(math.MaxBig256, big.NewInt(35)), big.NewInt(2))
	require.True(t, expectedChainID.Cmp(msg.ChainID()) == 0)
}

func TestGetTxFnSignatureInfo(t *testing.T) {
	chainID := big.NewInt(3)
	priv, _ := ethsecp256k1.GenerateKey()
	addr := ethcmn.BytesToAddress(priv.PubKey().Address().Bytes())
	msg := NewMsgEthereumTx(0, &addr, nil, 100000, nil, []byte("test"))
	err := msg.Sign(chainID, priv.ToECDSA())
	require.Nil(t, err)
}

func TestMsgEthereumTxGetter(t *testing.T) {
	priv, _ := ethsecp256k1.GenerateKey()
	addr := ethcmn.BytesToAddress(priv.PubKey().Address().Bytes())
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(100000)
	expectedFee := gasPrice * int64(gasLimit)
	expectCost := expectedFee + amount
	msg := NewMsgEthereumTx(0, &addr, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte("test"))

	require.Equal(t, gasLimit, msg.GetGas())
	require.True(t, big.NewInt(expectedFee).Cmp(msg.Fee()) == 0)
	require.True(t, big.NewInt(expectCost).Cmp(msg.Cost()) == 0)

	expectedV, expectedR, expectedS := big.NewInt(1), big.NewInt(2), big.NewInt(3)
	msg.Data.V, msg.Data.R, msg.Data.S = expectedV, expectedR, expectedS
	v, r, s := msg.RawSignatureValues()
	require.True(t, expectedV.Cmp(v) == 0)
	require.True(t, expectedR.Cmp(r) == 0)
	require.True(t, expectedS.Cmp(s) == 0)
}

func TestMsgEthereumTx_Amino(t *testing.T) {
	priv, _ := ethsecp256k1.GenerateKey()
	addr := ethcmn.BytesToAddress(priv.PubKey().Address().Bytes())
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(100000)
	msg := NewMsgEthereumTx(0, &addr, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte("test"))
	err := msg.Sign(big.NewInt(3), priv.ToECDSA())
	require.NoError(t, err)
	hash := ethcmn.BigToHash(big.NewInt(2))

	testCases := []*MsgEthereumTx{
		msg,
		{
			Data: TxData{
				AccountNonce: 2,
				Price:        big.NewInt(3),
				GasLimit:     1,
				Recipient:    &addr,
				Amount:       big.NewInt(4),
				Payload:      []byte("test"),
				V:            big.NewInt(5),
				R:            big.NewInt(6),
				S:            big.NewInt(7),
				Hash:         &hash,
			},
		},
		{
			Data: TxData{
				Price:     big.NewInt(math.MinInt64),
				Recipient: &ethcmn.Address{},
				Amount:    big.NewInt(math.MinInt64),
				Payload:   []byte{},
				V:         big.NewInt(math.MinInt64),
				R:         big.NewInt(math.MinInt64),
				S:         big.NewInt(math.MinInt64),
				Hash:      &ethcmn.Hash{},
			},
		},
		{
			Data: TxData{
				AccountNonce: math.MaxUint64,
				Price:        big.NewInt(math.MaxInt64),
				GasLimit:     math.MaxUint64,
				Amount:       big.NewInt(math.MaxInt64),
				V:            big.NewInt(math.MaxInt64),
				R:            big.NewInt(math.MaxInt64),
				S:            big.NewInt(math.MaxInt64),
			},
		},
	}

	for _, msg := range testCases {
		raw, err := ModuleCdc.MarshalBinaryBare(msg)
		require.NoError(t, err)

		var msg2 MsgEthereumTx
		err = ModuleCdc.UnmarshalBinaryBare(raw, &msg2)
		require.NoError(t, err)

		var msg3 MsgEthereumTx
		v, err := ModuleCdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(raw, &msg3)
		require.NoError(t, err)
		msg3 = *v.(*MsgEthereumTx)
		require.EqualValues(t, msg2, msg3)
	}
}

func BenchmarkMsgEthereumTxUnmarshal(b *testing.B) {
	cdc := ModuleCdc
	priv, _ := ethsecp256k1.GenerateKey()
	addr := ethcmn.BytesToAddress(priv.PubKey().Address().Bytes())
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(100000)
	msg := NewMsgEthereumTx(123456, &addr, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte("test"))
	_ = msg.Sign(big.NewInt(66), priv.ToECDSA())

	raw, _ := cdc.MarshalBinaryBare(msg)
	rlpRaw, err := rlp.EncodeToBytes(&msg)
	require.NoError(b, err)
	b.ResetTimer()

	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var msg2 MsgEthereumTx
			err := cdc.UnmarshalBinaryBare(raw, &msg2)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("unmarshaler", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var msg3 MsgEthereumTx
			v, err := cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(raw, &msg3)
			if err != nil {
				b.Fatal(err)
			}
			msg3 = v.(MsgEthereumTx)
		}
	})

	b.Run("rlp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var msg MsgEthereumTx
			err = rlp.DecodeBytes(rlpRaw, &msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestMarshalAndUnmarshalLogs(t *testing.T) {
	var cdc = codec.New()

	logs := []*ethtypes.Log{
		{
			Address: ethcmn.BytesToAddress([]byte{0x11}),
			TxHash:  ethcmn.HexToHash("0x01"),
			// May need to find workaround since Topics is required to unmarshal from JSON
			Topics:  []ethcmn.Hash{},
			Removed: true,
		},
		{Address: ethcmn.BytesToAddress([]byte{0x01, 0x11}), Topics: []ethcmn.Hash{}},
	}

	raw, err := codec.MarshalJSONIndent(cdc, logs)
	require.NoError(t, err)

	var logs2 []*ethtypes.Log
	err = cdc.UnmarshalJSON(raw, &logs2)
	require.NoError(t, err)

	require.Len(t, logs2, 2)
	require.Equal(t, logs[0].Address, logs2[0].Address)
	require.Equal(t, logs[0].TxHash, logs2[0].TxHash)
	require.True(t, logs[0].Removed)

	emptyLogs := []*ethtypes.Log{}

	raw, err = codec.MarshalJSONIndent(cdc, emptyLogs)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(raw, &logs2)
	require.NoError(t, err)
}

func TestMsgString(t *testing.T) {
	expectedUint64, expectedSDKAddr := uint64(1024), newSdkAddress()
	expectedPayload, err := hexutil.Decode("0x1234567890abcdef")
	require.NoError(t, err)
	expectedOutput := fmt.Sprintf("nonce=1024 gasPrice=1 gasLimit=1024 recipient=%s amount=1 data=0x1234567890abcdef from=%s",
		expectedSDKAddr, expectedSDKAddr)

	expectedHexAddr := ethcmn.BytesToAddress([]byte{0x01})
	expectedBigInt := big.NewInt(1024)
	expectedOutput = fmt.Sprintf("nonce=1024 price=1024 gasLimit=1024 recipient=%s amount=1024 data=0x1234567890abcdef v=0 r=0 s=0", expectedHexAddr.Hex())
	msgEthereumTx := NewMsgEthereumTx(expectedUint64, &expectedHexAddr, expectedBigInt, expectedUint64, expectedBigInt, expectedPayload)
	require.True(t, strings.EqualFold(msgEthereumTx.String(), expectedOutput))
}

func newProxyDecoder() *codec.CodecProxy {
	ModuleBasics := module.NewBasicManager(
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
	)
	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)
	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	return codecProxy
}
func TestMsgIBCTxValidate(t *testing.T) {
	tmtypes.UnittestOnlySetMilestoneVenus1Height(1)

	IBCRouterKey := "ibc"
	cpcdc := newProxyDecoder()
	marshaler := cpcdc.GetProtocMarshal()
	decode := ibctxdecode.IbcTxDecoder(marshaler)
	var err error
	txBytes1, err := hex.DecodeString("0a8d030a8a030a232f6962632e636f72652e636c69656e742e76312e4d7367437265617465436c69656e7412e2020aab010a2b2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436c69656e745374617465127c0a056962632d311204080110031a040880ac4d22040880df6e2a0308d80432003a05080110940342190a090801180120012a0100120c0a02000110211804200c300142190a090801180120012a0100120c0a02000110201801200130014a07757067726164654a1075706772616465644942435374617465500158011286010a2e2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436f6e73656e737573537461746512540a0c0892cbde930610a0ff9fe20212220a208acb6977f3cac564f6b015ff1de209e6c167e3454e6a754780e601efe340a5dd1a20cade35b27c5c32afead6cbed10d219c3903b8789b3fee9bf52b893efd6e2b8501a296578316a35657535716775376472686c737277346867326a766a3930707a766132373272347479763812720a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a210361469c236406f73459385bfe6d265ad8f293166b4661228d7bf6dd2f305236d912040a02080112200a1a0a0377656912133230343034393530303030303030303030303010e1a6081a40c951cde5885ab43d5e6c1ed88ef8adfd28311bfcba5461baa5bf4c9ad849e50837184dfde85ccb793f9859283553d3ef78113e5960aa353e885a9deb983e802a")
	txBytes2, err := hex.DecodeString("0af8080aeb070a232f6962632e636f72652e636c69656e742e76312e4d7367557064617465436c69656e7412c3070a0f30372d74656e6465726d696e742d301284070a262f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e48656164657212d9060ac9040a8c030a02080b12056962632d3118b907220c08d3dcde930610908ce5b3022a480a206241f13d50332dca5df961e0d1b0ab4c9dd36a8f05af2951ee23a0e1b1ff8913122408011220c2ed8f17294c68e683f31c0deaa1b34fe3966244a7ed52db9a716b9f9c200f703220d3f78b59cbb111d27452a4c0c71c614a6fada163ef175f8954f59007bc5d56df3a20f2473fff8995411fcde70619e3fa3a2e0a865e82edaae3c3ec072c8efaf10c014220c2c257672210f3023ae63cc03ac71b4517479b84ebee227719f06edd7b5aa60a4a20c2c257672210f3023ae63cc03ac71b4517479b84ebee227719f06edd7b5aa60a5220048091bc7ddc283f77bfbf91d73c44da58c3df8a9cbc867405d8b7f3daada22f5a208f988bd1fd57bc996fcba5b40814ee75d4dcc883d2f8eaa873b814ccdaa8acb66220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b8556a20e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85572143ff4a1ab7900d093bddea8d832c34f654207201a12b70108b9071a480a2025909a026179293b314c647573c4ee85b1a820b526e188cc42196e70b7d5622a122408011220ab82b5731b0b954a0fb947b897f522d51557c015632e040682edc607377111672268080212143ff4a1ab7900d093bddea8d832c34f654207201a1a0c08d4dcde93061090d0b8f2022240374206b5b047546cb08abd6994ddf6ec67a736940b286ea11627f20eb4b7668e8bc0ff55f6ab7a29e8fe23a98c672d23cfd926d363204408be729845f72fd1001280010a3e0a143ff4a1ab7900d093bddea8d832c34f654207201a12220a204d72a5c949a4140d889f39074a3302a07801134bb3a2f751ad7200cbb6da5a2e18a08d06123e0a143ff4a1ab7900d093bddea8d832c34f654207201a12220a204d72a5c949a4140d889f39074a3302a07801134bb3a2f751ad7200cbb6da5a2e18a08d061a05080110a4072280010a3e0a143ff4a1ab7900d093bddea8d832c34f654207201a12220a204d72a5c949a4140d889f39074a3302a07801134bb3a2f751ad7200cbb6da5a2e18a08d06123e0a143ff4a1ab7900d093bddea8d832c34f654207201a12220a204d72a5c949a4140d889f39074a3302a07801134bb3a2f751ad7200cbb6da5a2e18a08d061a29657831737061787a64376b37797a716467766b6e7568743671333870673530757530387666716135710a87010a2d2f6962632e636f72652e636f6e6e656374696f6e2e76312e4d7367436f6e6e656374696f6e4f70656e496e697412560a0f30372d74656e6465726d696e742d3012180a0f30372d74656e6465726d696e742d301a050a036962632a29657831737061787a64376b37797a716467766b6e75687436713338706735307575303876667161357112760a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a210311219a285ff5fd852d664a06279ce4cb60eb266be2cc2e24c33525895baade6612040a020801180112220a1c0a03776569121531363632323130303030303030303030303030303010cd920a1a403c793cbaa9d512f7e154511e1656a1cac1fff878d10b9b0a1d49596168bdad04786b9beb78fa1f58c33d580f1a09a345629211116dab2eda98276eb6edf05feb")
	txBytesArray := [][]byte{
		txBytes1, txBytes2,
	}
	expectedMsgType := []string{
		"/ibc.core.client.v1.MsgCreateClient",
		"/ibc.core.client.v1.MsgUpdateClient",
	}
	for i, txbytes := range txBytesArray {
		require.NoError(t, err)
		ibctx, err := decode(txbytes)
		require.NoError(t, err)
		require.NotNil(t, ibctx)
		require.Equal(t, ibctx.StdTx.Msgs[0].Route(), IBCRouterKey)
		require.Equal(t, ibctx.StdTx.Msgs[0].Type(), expectedMsgType[i])
		//tx validator
		require.NoError(t, ibctx.StdTx.Msgs[0].ValidateBasic())
	}
}

func TestMsgIbcTxMarshalSignBytes(t *testing.T) {
	chainID := "exchain-101"
	accnum := 1
	sequence := 0
	memo := "memo"
	authInfoBytes := []byte("authinfobytes")
	bodyBytes := []byte("bodyBytes")

	fee := authtypes.StdFee{
		Amount: []sdk.DecCoin{
			sdk.DecCoin{
				Denom:  "test",
				Amount: sdk.NewDecFromBigInt(big.NewInt(10)),
			},
		},
		Gas: 100000,
	}

	signBytes := authtypes.IbcDirectSignBytes(
		chainID,
		uint64(accnum),
		uint64(sequence),
		fee,
		nil,
		memo,
		authInfoBytes,
		bodyBytes,
	)

	expectedHexResult := "0A09626F64794279746573120D61757468696E666F62797465731A0B6578636861696E2D3130312001"

	require.Equal(t, expectedHexResult, fmt.Sprintf("%X", signBytes))

}

func BenchmarkEvmTxVerifySig(b *testing.B) {
	chainID := big.NewInt(3)
	priv1, _ := ethsecp256k1.GenerateKey()
	addr1 := ethcmn.BytesToAddress(priv1.PubKey().Address().Bytes())

	// require valid signature passes validation
	msg := NewMsgEthereumTx(0, &addr1, nil, 100000, nil, []byte("test"))
	_ = msg.Sign(chainID, priv1.ToECDSA())

	b.ResetTimer()

	b.Run("firstVerifySig", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := msg.firstVerifySig(chainID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestRlpPointerEncode(t *testing.T) {
	bz := make([]byte, 512)
	rand.Read(bz)

	h1 := rlpHash([]interface{}{bz})
	h2 := rlpHash([]interface{}{&bz})

	require.Equal(t, h1, h2)
}
