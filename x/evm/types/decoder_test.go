package types

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"math/big"
	"reflect"
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
		res, err = types.EthereumTxEncode(&expectedEthMsg)
	} else {
		res = cdc.MustMarshalBinaryLengthPrefixed(expectedEthMsg)
	}
	return
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}
func mustWtx(t *testing.T, cdc *codec.Codec, txbytes []byte) (wtx auth.WrappedTx) {
	decoder := TxDecoder(cdc)

	tx, err := decoder(txbytes, 2)
	require.NoError(t, err)

	var ok bool
	wtx, ok = tx.(auth.WrappedTx)
	require.Equal(t, ok, true)

	return
}

func TestAminoDecoder4EvmTx(t *testing.T) {

	cdc := makeCodec()
	decoder := TxDecoder(cdc)
	tmtypes.UnittestOnlySetMilestoneVenusHeight(1)
	defer tmtypes.UnittestOnlySetMilestoneVenusHeight(0)

	evmTxbytesByAmino, err := genEvmTxBytes(cdc, false)
	require.NoError(t, err)

	_, err = decoder(evmTxbytesByAmino, 2)
	t.Log(err)
	require.Error(t, err)

	_, err = ubDecoder(cdc, evmTxbytesByAmino, 2)
	t.Log(err)
	require.Error(t, err)

	_, err = ubruDecoder(cdc, evmTxbytesByAmino, 2)
	t.Log(err)
	require.Error(t, err)

	_, err = ubDecoder(cdc, evmTxbytesByAmino, 0)
	require.NoError(t, err)

	_, err = ubruDecoder(cdc, evmTxbytesByAmino, 0)
	require.NoError(t, err)
}

func TestWrappedTxDecoder(t *testing.T) {

	cdc := makeCodec()
	tmtypes.UnittestOnlySetMilestoneVenusHeight(1)
	defer tmtypes.UnittestOnlySetMilestoneVenusHeight(0)

	decoder := TxDecoder(cdc)

	evmTxbytesByRlp, err := genEvmTxBytes(cdc, true)
	require.NoError(t, err)

	var txBytesList [][]byte
	txBytesList = append(txBytesList, evmTxbytesByRlp)

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

		wtxBytes, err := types.EncodeWrappedTx(txbytes, info, int(sdk.EvmTxType))
		require.NoError(t, err)

		wtx, err := decoder(wtxBytes, 2)
		require.NoError(t, err)

		switch tx := wtx.(type) {
		case auth.WrappedTx:
			fmt.Printf("sdk.WrappedTx %+v\n", tx)
			break
		default:
			err = fmt.Errorf("received: %v", reflect.TypeOf(wtx).String())
		}
		require.NoError(t, err)
	}
}

func TestWrappedTxEncoder(t *testing.T) {

	cdc := makeCodec()
	tmtypes.UnittestOnlySetMilestoneVenusHeight(1)
	defer tmtypes.UnittestOnlySetMilestoneVenusHeight(0)

	evmTxbytesByRlp, err := genEvmTxBytes(cdc, true)
	require.NoError(t, err)

	info := &sdk.ExTxInfo{
		Metadata:  []byte("m1"),
		NodeKey:   []byte("n1"),
		Signature: []byte("s1"),
	}

	_, err = types.EncodeWrappedTx(evmTxbytesByRlp, info, int(sdk.WrappedTxType))
	require.Error(t, err)

	wtxBytes, err := types.EncodeWrappedTx(evmTxbytesByRlp, info, int(sdk.EvmTxType))
	require.NoError(t, err)

	wtx := mustWtx(t, cdc, wtxBytes)
	require.Equal(t, bytes.Compare(wtx.Metadata, info.Metadata), 0)
	require.Equal(t, bytes.Compare(wtx.NodeKey, info.NodeKey), 0)
	require.Equal(t, bytes.Compare(wtx.Signature, info.Signature), 0)

	info2 := &sdk.ExTxInfo{
		Metadata:  []byte("m2"),
		NodeKey:   []byte("n2"),
		Signature: []byte("s2"),
	}

	wtxBytes, err = types.EncodeWrappedTx(wtxBytes, info2, int(sdk.WrappedTxType))
	require.NoError(t, err)

	wtx = mustWtx(t, cdc, wtxBytes)
	require.Equal(t, bytes.Compare(wtx.Metadata, info2.Metadata), 0)
	require.Equal(t, bytes.Compare(wtx.NodeKey, info2.NodeKey), 0)
	require.Equal(t, bytes.Compare(wtx.Signature, info2.Signature), 0)

	// todo
	//wtxBytes, err = types.EncodeWrappedTx(wtxBytes, info2, int(types.EvmTxType))
	//require.Error(t, err)
}
