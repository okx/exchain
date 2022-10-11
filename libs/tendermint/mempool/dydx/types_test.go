package dydx

import (
	"encoding/hex"
	"github.com/okex/exchain/libs/dydx/contracts"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

const (
	orderHex       = "000000000000000000000000000000000000000000000000000000000000000173646861000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000004000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb9226600000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c80000000000000000000000000000000000000000000000000000000000000005"
	signedOrderHex = "000000000000000000000000000000000000000000000000000000000000004073646861617364617364000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000000173646861000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000004000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb9226600000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c80000000000000000000000000000000000000000000000000000000000000005"
	orderSigHex    = "7364686161736461736400000000000000000000000000000000000000000000"
	flagsHex       = "7364686100000000000000000000000000000000000000000000000000000000"
	makerHex       = "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	takerHex       = "70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

func TestDecodeSignedMsg(t *testing.T) {
	signedMsgBytes, err := hex.DecodeString(signedOrderHex)
	require.NoError(t, err)
	var so SignedOrder
	err = so.DecodeFrom(signedMsgBytes)
	require.NoError(t, err)
	require.Equal(t, orderHex, hex.EncodeToString(so.Msg))
	require.Equal(t, orderSigHex, hex.EncodeToString(so.Sig[:]))
}

func TestEncode(t *testing.T) {
	a := int32(17)
	data, err := callTypeABI.Encode(a)
	require.NoError(t, err)
	require.Equal(t, byte(a), data[len(data)-1])
}

func TestDecodeOrder(t *testing.T) {
	flagsBytes, err := hex.DecodeString(flagsHex)
	require.NoError(t, err)
	makerBytes, err := hex.DecodeString(makerHex)
	require.NoError(t, err)
	takerBytes, err := hex.DecodeString(takerHex)
	require.NoError(t, err)
	odr := P1Order{
		CallType: 1,
		P1OrdersOrder: contracts.P1OrdersOrder{
			Amount:       big.NewInt(1),
			LimitPrice:   big.NewInt(2),
			TriggerPrice: big.NewInt(3),
			LimitFee:     big.NewInt(4),
			Expiration:   big.NewInt(5),
		},
	}
	copy(odr.Flags[:], flagsBytes)
	copy(odr.Maker[:], makerBytes)
	copy(odr.Taker[:], takerBytes)
	orderBytes, err := odr.Encode()
	require.NoError(t, err)
	require.Equal(t, orderHex, hex.EncodeToString(orderBytes))

	var odr2 P1Order
	err = odr2.DecodeFrom(orderBytes)
	require.NoError(t, err)
	require.Equal(t, odr, odr2)
}

func TestDecodeRLP(t *testing.T) {
	data := "f8ab819f843b9aca00832dc6c0945d64795f3f815924e607c7e9651e89db4dbddb6280b844a9059cbb00000000000000000000000033c866e121fa09a23a7dbecb87ad9c394d3d452300000000000000000000000000000000000000000000000000000002540be40081aaa08838daee659574adbea5efb9c36c3901a8d33275122403d10eec9c1bab461be5a0446b28b2014bf7b490e2297f19202fd6290b0e82657713df9661ee21b78a647e"
	txBytes, err := hex.DecodeString(data)
	require.NoError(t, err)
	var tx TxData
	err = rlp.DecodeBytes(txBytes, &tx)
	require.NoError(t, err)
	require.Equal(t, "0x5D64795f3f815924E607C7e9651e89Db4Dbddb62", tx.Recipient.String())
}

func TestDecodeParallel(t *testing.T) {
	var wg sync.WaitGroup
	f := func() {
		defer wg.Done()
		signedMsgBytes, err := hex.DecodeString(signedOrderHex)
		require.NoError(t, err)
		var so SignedOrder
		err = so.DecodeFrom(signedMsgBytes)
		require.NoError(t, err)
		require.Equal(t, orderHex, hex.EncodeToString(so.Msg))
		require.Equal(t, orderSigHex, hex.EncodeToString(so.Sig[:]))
	}

	total := 10000
	wg.Add(total)
	for i := 0; i < total; i++ {
		f()
	}
	wg.Wait()
}

func TestHash(t *testing.T) {
	hashHex := "0x8e67b6484476c8e2a168e37bed7be6212d5aaedc08869f8d6083fffdee2eb3ea"
	orderBytes, err := hex.DecodeString(orderHex)
	require.NoError(t, err)
	var odr P1Order
	err = odr.DecodeFrom(orderBytes)
	require.NoError(t, err)
	require.Equal(t, hashHex, odr.Hash().String())
	odr.CallType += 1
	require.Equal(t, hashHex, odr.Hash().String())
}
