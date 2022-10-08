package dydx

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	testnetChainID := big.NewInt(65)
	// ethRpcUrl := "https://exchaintestrpc.okex.org"
	ethWsUrl := "wss://exchaintestws.okex.org:8443"
	fromBlockNum := big.NewInt(14704890)
	endBlockNum := big.NewInt(14704893)
	privKey := "e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d"

	client, err := NewDydxClient(testnetChainID, ethWsUrl, fromBlockNum, privKey,
		"0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
		"0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619")
	require.NoError(t, err)

	endBlock := endBlockNum.Uint64()
	iter, err := client.contracts.PerpetualV1.FilterLogTrade(&bind.FilterOpts{
		Start:   fromBlockNum.Uint64(),
		End:     &endBlock,
		Context: context.Background(),
	}, nil, nil)
	require.NoError(t, err)
	for iter.Next() {
		t.Logf("LogTrade: %+v", iter.Event)
	}
	_ = iter.Close()

	client.Stop()
}
