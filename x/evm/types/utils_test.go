package types

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
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

	enc, err := EncodeResultData(data)
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
			BlockNumber: 18,
			TxHash:      ethcmn.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			TxIndex:     0,
			BlockHash:   ethcmn.HexToHash("0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF"),
			Index:       0,
			Removed:     true,
		},
	}
	cdc := codec.New()
	for _, test := range tests {
		bz, err := cdc.MarshalBinaryBare(test)
		require.NoError(t, err)

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
		{
			ContractAddress: addr,
			Bloom:           bloom,
			Logs: []*ethtypes.Log{
				{
					Data:        []byte{1, 2, 3, 4},
					BlockNumber: 17,
					Index:       10,
				}},
			Ret:    ret,
			TxHash: ethcmn.HexToHash("0x00"),
		},
		{},
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

		actual, err := data.MarshalToAmino()
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)
		t.Log(fmt.Sprintf("%d pass\n", i))
	}
}
