package types

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	okexchaincodec "github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	ibctransfer "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/require"
)

func TestEvmDataEncoding(t *testing.T) {
	addr := ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49")
	bloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	ret := []byte{0x5, 0x8}

	data := ResultData{
		ContractAddress: addr,
		Bloom:           bloom,
		Logs: []*ethtypes.Log{{
			Data:        []byte{1, 2, 3, 4},
			BlockNumber: 17,
		}},
		Ret: ret,
	}

	enc, err := EncodeResultData(&data)
	require.NoError(t, err)

	res, err := DecodeResultData(enc)
	require.NoError(t, err)
	require.Equal(t, addr, res.ContractAddress)
	require.Equal(t, bloom, res.Bloom)
	require.Equal(t, data.Logs, res.Logs)
	require.Equal(t, ret, res.Ret)

	// error check
	_, err = DecodeResultData(enc[1:])
	require.Error(t, err)
}

func TestValidateSigner(t *testing.T) {
	const digest = "default digest"
	digestHash := crypto.Keccak256([]byte(digest))
	priv, err := crypto.GenerateKey()
	require.NotNil(t, priv)
	require.NoError(t, err)

	ethAddr := crypto.PubkeyToAddress(priv.PublicKey)
	require.NoError(t, err)

	sig, err := crypto.Sign(digestHash, priv)
	require.NoError(t, err)

	err = ValidateSigner(digestHash, sig, ethAddr)
	require.NoError(t, err)

	// different eth address
	otherEthAddr := ethcmn.BytesToAddress([]byte{1})
	err = ValidateSigner(digestHash, sig, otherEthAddr)
	require.Error(t, err)

	// invalid digestHash
	err = ValidateSigner(digestHash[1:], sig, otherEthAddr)
	require.Error(t, err)
}

func TestResultData_String(t *testing.T) {
	const expectedResultDataStr = `ResultData:
	ContractAddress: 0x5dE8a020088a2D6d0a23c204FFbeD02790466B49
	Bloom: 259
	Ret: [5 8]
	TxHash: 0x0000000000000000000000000000000000000000000000000000000000000000	
	Logs: 
		{0x0000000000000000000000000000000000000000 [] [1 2 3 4] 17 0x0000000000000000000000000000000000000000000000000000000000000000 0 0x0000000000000000000000000000000000000000000000000000000000000000 0 false}
 		{0x0000000000000000000000000000000000000000 [] [5 6 7 8] 18 0x0000000000000000000000000000000000000000000000000000000000000000 0 0x0000000000000000000000000000000000000000000000000000000000000000 0 false}`
	addr := ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49")
	bloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	ret := []byte{0x5, 0x8}

	data := ResultData{
		ContractAddress: addr,
		Bloom:           bloom,
		Logs: []*ethtypes.Log{
			{
				Data:        []byte{1, 2, 3, 4},
				BlockNumber: 17,
			},
			{
				Data:        []byte{5, 6, 7, 8},
				BlockNumber: 18,
			}},
		Ret: ret,
	}

	require.True(t, strings.EqualFold(expectedResultDataStr, data.String()))
}

func newProxyDecoder() *codec.CompoundCodec {
	ModuleBasics := module.NewBasicManager(
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
	)
	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)
	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	return &codec.CompoundCodec{
		cdc,
		codecProxy,
	}
}

func TestTxDecoderForIbcSuccess(t *testing.T) {
	cpcdc := newProxyDecoder()
	var IBCRouterKey = "ibc"

	txDecoder := TxDecoder(cpcdc)

	txBytes1, _ := hex.DecodeString("0a8c030a89030a232f6962632e636f72652e636c69656e742e76312e4d7367437265617465436c69656e7412e1020aaa010a2b2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436c69656e745374617465127b0a056962632d311204080110031a040880ac4d22040880df6e2a0308d80432003a040801102142190a090801180120012a0100120c0a02000110211804200c300142190a090801180120012a0100120c0a02000110201801200130014a07757067726164654a1075706772616465644942435374617465500158011286010a2e2f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e436f6e73656e737573537461746512540a0c08a2eb9992061080a4b1b60312220a206c4821c63adcb62c57effb912d6d568f8971456acf6630c5b3677c8d3cd999321a203ee89cb19deef923d4ffea8d44980336746dd0125d7372f9a78562a8aaf0fb821a296578316865756b7866766d74617968306a7a746468706e73756468677a666d3364783772357261753412370a290a210a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912040a020801120a0a080a036f6b741201301a00")
	txBytes2, _ := hex.DecodeString("0af4080ae7070a232f6962632e636f72652e636c69656e742e76312e4d7367557064617465436c69656e7412bf070a0f30372d74656e6465726d696e742d301280070a262f6962632e6c69676874636c69656e74732e74656e6465726d696e742e76312e48656164657212d5060ac6040a8a030a02080b12056962632d311823220b08a5eb99920610c0f19c6e2a480a20629665f068bb27c518e135e3b830759d7dbbb4b57bc77671de22b6baae690e911224080112203729aa9acb11d1ca83dff915dc0f98d7ae71958422c239c3b1b50f0cc6a2f90c322082e4414794afbce9acad186181a52d26960917b5b042713c7a6207a91d31720f3a2017ed596fab8c5fa97aeec4e927b6099ba7b1b70e09baa9d10e819f749b7f35ec42203ee89cb19deef923d4ffea8d44980336746dd0125d7372f9a78562a8aaf0fb824a203ee89cb19deef923d4ffea8d44980336746dd0125d7372f9a78562a8aaf0fb825220048091bc7ddc283f77bfbf91d73c44da58c3df8a9cbc867405d8b7f3daada22f5a20748be931993627eb285a0eae322aa3595a4337578abd499ec0af1a83f38307ea6220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b8556a20e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85572148bbcac9722c1a67a9c7517ad320f9c547bdf1a7512b60108231a480a20f482ea7b50914243a41fbcbee704ee17b7b268cee54519fd94ca66b1cc0f59241224080112208008d042c3261bc12e3007cfcc32e0cd9246c36a1f82fcdd4c1c3a0ff714437f2268080212148bbcac9722c1a67a9c7517ad320f9c547bdf1a751a0c08a6eb9992061090f1c9c20122404b2117e0b4d41912726a03a9179425d61c82114c80b5fbdf0542bb7bfde5278843e148bce5246e957877ecba92df1b81bd5d46e5c8e03bb5bed2f91fc0c5020b1280010a3e0a148bbcac9722c1a67a9c7517ad320f9c547bdf1a7512220a208ad3e4b95873ec49c803d52bddff903b50f72faf5878dc26e55e5f62dfea32fa18a08d06123e0a148bbcac9722c1a67a9c7517ad320f9c547bdf1a7512220a208ad3e4b95873ec49c803d52bddff903b50f72faf5878dc26e55e5f62dfea32fa18a08d061a04080110212280010a3e0a148bbcac9722c1a67a9c7517ad320f9c547bdf1a7512220a208ad3e4b95873ec49c803d52bddff903b50f72faf5878dc26e55e5f62dfea32fa18a08d06123e0a148bbcac9722c1a67a9c7517ad320f9c547bdf1a7512220a208ad3e4b95873ec49c803d52bddff903b50f72faf5878dc26e55e5f62dfea32fa18a08d061a296578316865756b7866766d74617968306a7a746468706e73756468677a666d336478377235726175340a87010a2d2f6962632e636f72652e636f6e6e656374696f6e2e76312e4d7367436f6e6e656374696f6e4f70656e496e697412560a0f30372d74656e6465726d696e742d3012180a0f30372d74656e6465726d696e742d301a050a036962632a296578316865756b7866766d74617968306a7a746468706e73756468677a666d3364783772357261753412390a2b0a210a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912040a0208011801120a0a080a036f6b741201301a00")
	txBytesArray := [][]byte{
		txBytes1, txBytes2,
	}
	expectedMsgType := []string{
		"create_client",
		"update_client",
	}
	for i, txbytes := range txBytesArray {
		tx, err := txDecoder(txbytes)
		ibcTx, _ := tx.(*authtypes.IbcTx)

		require.NoError(t, err)
		require.NotNil(t, ibcTx)
		require.Equal(t, ibcTx.StdTx.Msgs[0].Route(), IBCRouterKey)
		require.Equal(t, ibcTx.StdTx.Msgs[0].Type(), expectedMsgType[i])
		require.NoError(t, ibcTx.StdTx.Msgs[0].ValidateBasic())
	}
}

func TestTxDecoderForIbcFailure(t *testing.T) {
	cpcdc := newProxyDecoder()
	txDecoder := TxDecoder(cpcdc)
	txBytes1, _ := hex.DecodeString("123456")
	_, err := txDecoder(txBytes1)

	require.Error(t, err)
}

func TestTxDecoder(t *testing.T) {
	expectUint64, expectedBigInt, expectedBytes := uint64(1024), big.NewInt(1024), []byte("default payload")
	expectedEthAddr := ethcmn.BytesToAddress([]byte("test_address"))
	expectedEthMsg := NewMsgEthereumTx(expectUint64, &expectedEthAddr, expectedBigInt, expectUint64, expectedBigInt, expectedBytes)

	// register codec
	cdc := codec.New()
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)
	RegisterCodec(cdc)

	txbytes := cdc.MustMarshalBinaryLengthPrefixed(expectedEthMsg)
	txDecoder := TxDecoder(cdc)
	tx, err := txDecoder(txbytes)
	require.NoError(t, err)

	msgs := tx.GetMsgs()
	require.Equal(t, 1, len(msgs))
	require.NoError(t, msgs[0].ValidateBasic())
	require.True(t, strings.EqualFold(expectedEthMsg.Route(), msgs[0].Route()))
	require.True(t, strings.EqualFold(expectedEthMsg.Type(), msgs[0].Type()))

	require.NoError(t, tx.ValidateBasic())

	// error check
	_, err = txDecoder([]byte{})
	require.Error(t, err)

	_, err = txDecoder(txbytes[1:])
	require.Error(t, err)

	oldHeight := types.GetMilestoneVenusHeight()
	defer types.UnittestOnlySetMilestoneVenusHeight(oldHeight)
	rlpBytes, err := rlp.EncodeToBytes(&expectedEthMsg)
	require.Nil(t, err)

	for _, c := range []struct {
		curHeight          int64
		venusHeight        int64
		enableAminoDecoder bool
		enableRLPDecoder   bool
	}{
		{999, 0, true, false},
		{999, 1000, true, false},
		{1000, 1000, false, true},
		{1500, 1000, false, true},
	} {
		types.UnittestOnlySetMilestoneVenusHeight(c.venusHeight)
		_, err = TxDecoder(cdc)(txbytes, c.curHeight)
		require.Equal(t, c.enableAminoDecoder, err == nil)
		_, err = TxDecoder(cdc)(rlpBytes, c.curHeight)
		require.Equal(t, c.enableRLPDecoder, err == nil)

		// use global height when height is not pass through parameters.
		global.SetGlobalHeight(c.curHeight)
		_, err = TxDecoder(cdc)(txbytes)
		require.Equal(t, c.enableAminoDecoder, err == nil)
		_, err = TxDecoder(cdc)(rlpBytes)
		require.Equal(t, c.enableRLPDecoder, err == nil)
	}
}

func TestEthLogAmino(t *testing.T) {
	tests := []ethtypes.Log{
		{},
		{Topics: []ethcmn.Hash{}, Data: []byte{}},
		{
			Address: ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49"),
			Topics: []ethcmn.Hash{
				ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				ethcmn.HexToHash("0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"),
			},
			Data:        []byte{1, 2, 3, 4},
			BlockNumber: 17,
			TxHash:      ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			TxIndex:     123456,
			BlockHash:   ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			Index:       543121,
			Removed:     false,
		},
		{
			Address: ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49"),
			Topics: []ethcmn.Hash{
				ethcmn.HexToHash("0x00000000FF0000000000000000000AC0000000000000EF000000000000000000"),
				ethcmn.HexToHash("0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"),
			},
			Data:        []byte{5, 6, 7, 8},
			BlockNumber: math.MaxUint64,
			TxHash:      ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			TxIndex:     math.MaxUint,
			BlockHash:   ethcmn.HexToHash("0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"),
			Index:       math.MaxUint,
			Removed:     true,
		},
	}
	cdc := codec.New()
	for _, test := range tests {
		bz, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

		bz2, err := MarshalEthLogToAmino(&test)
		require.NoError(t, err)
		require.EqualValues(t, bz, bz2)

		var expect ethtypes.Log
		err = cdc.UnmarshalBinaryBare(bz, &expect)
		require.NoError(t, err)

		actual, err := UnmarshalEthLogFromAmino(bz)
		require.NoError(t, err)
		require.EqualValues(t, expect, *actual)
	}
}

func TestResultDataAmino(t *testing.T) {
	addr := ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49")
	bloom := ethtypes.BytesToBloom([]byte{0x1, 0x3, 0x5, 0x7})
	ret := []byte{0x5, 0x8}

	cdc := codec.New()
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)
	RegisterCodec(cdc)

	testDataSet := []ResultData{
		{},
		{Logs: []*ethtypes.Log{}, Ret: []byte{}},
		{
			ContractAddress: addr,
			Bloom:           bloom,
			Logs: []*ethtypes.Log{
				{
					Data:        []byte{1, 2, 3, 4},
					BlockNumber: 17,
					Index:       10,
				},
				{
					Data:        []byte{1, 2, 3, 4},
					BlockNumber: 17,
					Index:       10,
				},
				{
					Data:        []byte{1, 2, 3, 4},
					BlockNumber: 17,
					Index:       10,
				},
				nil,
			},
			Ret:    ret,
			TxHash: ethcmn.HexToHash("0x00"),
		},
		{
			ContractAddress: addr,
			Bloom:           bloom,
			Logs: []*ethtypes.Log{
				nil,
				{
					Removed: true,
				},
			},
			Ret:    ret,
			TxHash: ethcmn.HexToHash("0x00"),
		},
	}

	for i, data := range testDataSet {
		expect, err := cdc.MarshalBinaryBare(data)
		require.NoError(t, err)

		actual, err := data.MarshalToAmino(cdc)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)
		t.Log(fmt.Sprintf("%d pass\n", i))

		var expectRd ResultData
		err = cdc.UnmarshalBinaryBare(expect, &expectRd)
		require.NoError(t, err)
		var actualRd ResultData
		err = actualRd.UnmarshalFromAmino(cdc, expect)
		require.NoError(t, err)
		require.EqualValues(t, expectRd, actualRd)

		encoded, err := EncodeResultData(&data)
		require.NoError(t, err)
		decodedRd, err := DecodeResultData(encoded)
		require.NoError(t, err)
		require.EqualValues(t, expectRd, decodedRd)
	}
}

func BenchmarkDecodeResultData(b *testing.B) {
	addr := ethcmn.HexToAddress("0x5dE8a020088a2D6d0a23c204FFbeD02790466B49")
	bloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	ret := []byte{0x5, 0x8}

	data := ResultData{
		ContractAddress: addr,
		Bloom:           bloom,
		Logs: []*ethtypes.Log{{
			Data:        []byte{1, 2, 3, 4},
			BlockNumber: 17,
		}},
		Ret:    ret,
		TxHash: ethcmn.BigToHash(big.NewInt(10)),
	}

	enc, err := EncodeResultData(&data)
	require.NoError(b, err)
	b.ResetTimer()
	b.Run("amino", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var rd ResultData
			err = ModuleCdc.UnmarshalBinaryLengthPrefixed(enc, &rd)
			if err != nil {
				panic("err should be nil")
			}
		}
	})
	b.Run("unmarshaler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = DecodeResultData(enc)
			if err != nil {
				panic("err should be nil")
			}
		}
	})
}

func TestEthStringer(t *testing.T) {
	max := 10
	wg := &sync.WaitGroup{}
	wg.Add(max)
	for i := 0; i < max; i++ {
		go func() {
			addr := GenerateEthAddress()
			h := addr.Hash()
			require.Equal(t, addr.String(), EthAddressStringer(addr).String())
			require.Equal(t, h.String(), EthHashStringer(h).String())
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkEthAddressStringer(b *testing.B) {
	addr := GenerateEthAddress()
	b.ResetTimer()
	b.Run("eth", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = addr.String()
		}
	})
	b.Run("oec stringer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = EthAddressStringer(addr).String()
		}
	})
}

func BenchmarkEthHashStringer(b *testing.B) {
	addr := GenerateEthAddress()
	h := addr.Hash()
	b.ResetTimer()
	b.Run("eth", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = h.String()
		}
	})
	b.Run("oec stringer", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = EthHashStringer(h).String()
		}
	})
}
