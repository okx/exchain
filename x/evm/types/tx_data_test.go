package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

func TestMarshalAndUnmarshalData(t *testing.T) {
	addr := GenerateEthAddress()
	hash := ethcmn.BigToHash(big.NewInt(2))

	txData := TxData{
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
	}

	bz, err := txData.MarshalAmino()
	require.NoError(t, err)
	require.NotNil(t, bz)

	var txData2 TxData
	err = txData2.UnmarshalAmino(bz)
	require.NoError(t, err)

	require.Equal(t, txData, txData2)

	// check error
	err = txData2.UnmarshalAmino(bz[1:])
	require.Error(t, err)
}

func TestMsgEthereumTxAmino(t *testing.T) {
	addr := GenerateEthAddress()
	msg := NewMsgEthereumTx(5, &addr, big.NewInt(1), 100000, big.NewInt(3), []byte("test"))

	msg.Data.V = big.NewInt(1)
	msg.Data.R = big.NewInt(2)
	msg.Data.S = big.NewInt(3)

	raw, err := ModuleCdc.MarshalBinaryBare(msg)
	require.NoError(t, err)

	var msg2 MsgEthereumTx

	err = ModuleCdc.UnmarshalBinaryBare(raw, &msg2)
	require.NoError(t, err)
	require.Equal(t, msg, msg2)
}

func TestTxData_String(t *testing.T) {
	const expectedStrWithoutRecipient = "nonce=2 price=3 gasLimit=1 recipient=nil amount=4 data=0x1234567890abcdef v=5 r=6 s=7"
	payload, err := hexutil.Decode("0x1234567890abcdef")
	require.NoError(t, err)
	txData := TxData{
		AccountNonce: 2,
		Price:        big.NewInt(3),
		GasLimit:     1,
		Amount:       big.NewInt(4),
		Payload:      payload,
		V:            big.NewInt(5),
		R:            big.NewInt(6),
		S:            big.NewInt(7),
	}

	require.True(t, strings.EqualFold(expectedStrWithoutRecipient, txData.String()))

	// add recipient
	const expectedStrWithRecipient = "nonce=2 price=3 gasLimit=1 recipient=0x0000000000000000000000000000000000000000 amount=4 data=0x1234567890abcdef v=5 r=6 s=7"
	expectedEthAddr := ethcmn.HexToAddress("0x0000000000000000000000000000000000000000")
	txData.Recipient = &expectedEthAddr
	require.True(t, strings.EqualFold(expectedStrWithRecipient, txData.String()))
}
