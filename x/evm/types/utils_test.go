package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
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
