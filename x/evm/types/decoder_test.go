package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/okex/exchain/x/staking/types"
	"math/big"
	"reflect"
	"strings"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)


func genEvmTxBytes(cdc *codec.Codec, rlp bool) (res []byte, err error) {
	expectUint64, expectedBigInt, expectedBytes := uint64(1024), big.NewInt(1024), []byte("default payload")
	expectedEthAddr := ethcmn.BytesToAddress([]byte("test_address"))
	expectedEthMsg := NewMsgEthereumTx(expectUint64, &expectedEthAddr, expectedBigInt, expectUint64, expectedBigInt, expectedBytes)
	if rlp {
		res, err = types.EncodeEvmTx(&expectedEthMsg)
	} else {
		res = cdc.MustMarshalBinaryLengthPrefixed(expectedEthMsg)
	}
	return
}


func genTxBytes(cdc *codec.Codec) (res []byte, err error) {

	msg := stakingtypes.MsgEditValidator{
		Description: stakingtypes.Description{
			"1",
			"12",
			"3",
			"4",
		},
	}
	stakingtypes.RegisterCodec(cdc)
	res, err = cdc.MarshalBinaryLengthPrefixed(msg)
	return
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}
func TestCheckedTxDecoder(t *testing.T) {

	cdc := codec.New()
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)
	RegisterCodec(cdc)
	decoder := TxDecoder(cdc)

	//evmTxbytesByRlp, err := genEvmTxBytes(cdc, true)
	//require.NoError(t, err)

	evmTxbytesByAmino, err := genEvmTxBytes(cdc, false)
	require.NoError(t, err)

	//cmTxbytesByAmino, err := genTxBytes(cdc)
	//require.NoError(t, err)

	var txBytesList [][]byte
	//txBytesList = append(txBytesList, cmTxbytesByAmino)
	//txBytesList = append(txBytesList, evmTxbytesByRlp)
	txBytesList = append(txBytesList, evmTxbytesByAmino)

	for _, txbytes := range txBytesList {
		evmTx, err := decoder(txbytes, 2)
		require.NoError(t, err)

		switch tx := evmTx.(type) {
		case MsgEthereumTx:
			fmt.Printf("MsgEthereumTx %+v\n", tx)
		default:
			err = fmt.Errorf("received: %v", reflect.TypeOf(evmTx).String())
		}
		require.NoError(t, err)

		info := &sdk.ExTxInfo{
			Metadata:  []byte("m1"),
			NodeKey:   []byte("n1"),
			Signature: []byte("s1"),
		}

		chkTxBytes, err := types.EncodeWrappedTx(txbytes, info, false)
		require.NoError(t, err)

		chkTx, err := decoder(chkTxBytes, 2)
		require.NoError(t, err)

		switch tx := chkTx.(type) {
		case auth.WrappedTx:
			fmt.Printf("sdk.CheckedTx %+v\n", tx)
			break
		default:
			err = fmt.Errorf("received: %v", reflect.TypeOf(chkTx).String())
		}
		require.NoError(t, err)
	}
	//txDecoder := TxDecoder(cdc)
	//tx, err := txDecoder(txbytes)
	//require.NoError(t, err)

	//msgs := tx.GetMsgs()
	//require.Equal(t, 1, len(msgs))
	//require.NoError(t, msgs[0].ValidateBasic())
	//require.True(t, strings.EqualFold(expectedEthMsg.Route(), msgs[0].Route()))
	//require.True(t, strings.EqualFold(expectedEthMsg.Type(), msgs[0].Type()))
	//
	//require.NoError(t, tx.ValidateBasic())
	//
	//// error check
	//_, err = txDecoder([]byte{})
	//require.Error(t, err)
	//
	//_, err = txDecoder(txbytes[1:])
	//require.Error(t, err)
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
}
