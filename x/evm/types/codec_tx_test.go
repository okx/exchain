package types

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

type encoder interface {
	encodeTx(tx sdk.Tx) ([]byte, error)
	decodeTx(b []byte, tx interface{}) error
	name() string
}

type encodeMode uint64

const (
	rawAminoEnc encodeMode = iota
	rlpEnc
	exAminoEnc
)

var (
	rawEthMsgName = "raw/eth_tx"
	exEthMsgName  = "exchain/eth_tx"
)

func generateTestTx(n int) MsgEthereumTx {
	nonce := uint64(0)
	to := common.HexToAddress("0x5E7BA03cf5394c3731242164b294a968d9D937F2")
	amount := new(big.Int).SetUint64(100)
	gasLimit := uint64(21000)
	gasPrice, _ := new(big.Int).SetString("1000000", 10)
	data := strings.Repeat("1234567890", n)
	return NewMsgEthereumTx(nonce, &to, amount, gasLimit, gasPrice, []byte(data))
}

func newTestEncoder(mode encodeMode) encoder {
	switch mode {
	case rawAminoEnc:
		return newRawAminoEncoder()
	case rlpEnc:
		return newRlpEncoder()
	case exAminoEnc:
		return newExAminoEncoder()
	default:
	}
	panic("unknow encoder")
}

type rawAminoEncoder struct {
	cdc *amino.Codec
}

func newRawAminoEncoder() *rawAminoEncoder {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)
	cdc.RegisterConcrete(&MsgEthereumTx{}, rawEthMsgName, nil)

	return &rawAminoEncoder{cdc: cdc}
}

func (re *rawAminoEncoder) encodeTx(tx sdk.Tx) ([]byte, error) {
	return re.cdc.MarshalBinaryLengthPrefixed(tx)
}

func (re *rawAminoEncoder) decodeTx(b []byte, tx interface{}) error {
	return re.cdc.UnmarshalBinaryLengthPrefixed(b, tx)
}
func (re *rawAminoEncoder) name() string { return "go-amino" }

type rlpEncoder struct{}

func newRlpEncoder() *rlpEncoder {
	return &rlpEncoder{}
}

func (re *rlpEncoder) encodeTx(tx sdk.Tx) ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

func (re *rlpEncoder) decodeTx(b []byte, tx interface{}) error {
	return rlp.DecodeBytes(b, tx)
}
func (re *rlpEncoder) name() string { return "rlp " }

type exAminoEncoder struct {
	cdc *amino.Codec
}

func newExAminoEncoder() *exAminoEncoder {
	cdc := codec.New()
	cdc.RegisterConcrete(MsgEthereumTx{}, exEthMsgName, nil)
	cdc.RegisterConcreteUnmarshaller(exEthMsgName, func(_ *amino.Codec, bytes []byte) (interface{}, int, error) {
		var msg MsgEthereumTx
		err := msg.UnmarshalFromAmino(bytes)
		if err != nil {
			return nil, 0, err
		}
		return msg, len(bytes), nil
	})

	return &exAminoEncoder{cdc: cdc}
}

func (ee *exAminoEncoder) encodeTx(tx sdk.Tx) ([]byte, error) {
	return ee.cdc.MarshalBinaryLengthPrefixed(tx)
}

func (ee *exAminoEncoder) decodeTx(b []byte, tx interface{}) error {
	_, err := ee.cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(b, tx)
	return err
}
func (ee *exAminoEncoder) name() string { return "exchain-amino" }

func TestEncoder(t *testing.T) {
	testEncoder(t, newTestEncoder(rawAminoEnc)) //test go-amino
	testEncoder(t, newTestEncoder(rlpEnc))      //test ethereum-rlp
	testEncoder(t, newTestEncoder(exAminoEnc))  //test exchain-amino
}
func testEncoder(t *testing.T, enc encoder) {
	// encode
	tx := generateTestTx(1)
	data, err := enc.encodeTx(&tx)
	require.NoError(t, err, enc.name())

	// decode
	evmMsg := new(MsgEthereumTx)
	err = enc.decodeTx(data, evmMsg)
	require.NoError(t, err, enc.name())
}
func BenchmarkRawAminoEncodeTx(b *testing.B) { benchmarkEncodeTx(b, newTestEncoder(rawAminoEnc)) }
func BenchmarkRlpEncodeTx(b *testing.B)      { benchmarkEncodeTx(b, newTestEncoder(rlpEnc)) }
func BenchmarkExAminoEncodeTx(b *testing.B)  { benchmarkEncodeTx(b, newTestEncoder(exAminoEnc)) }

func benchmarkEncodeTx(b *testing.B, enc encoder) {
	tx := generateTestTx(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.encodeTx(&tx)
	}
}

func BenchmarkRawAminoDecodeTx(b *testing.B) { benchmarkDecodeTx(b, newTestEncoder(rawAminoEnc)) }
func BenchmarkRlpDecodeTx(b *testing.B)      { benchmarkDecodeTx(b, newTestEncoder(rlpEnc)) }
func BenchmarkExDecodeTx(b *testing.B)       { benchmarkDecodeTx(b, newTestEncoder(exAminoEnc)) }

func benchmarkDecodeTx(b *testing.B, enc encoder) {
	tx := generateTestTx(1)
	data, _ := enc.encodeTx(&tx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		evmMsg := new(MsgEthereumTx)
		b.StartTimer()
		enc.decodeTx(data, evmMsg)
	}
}
