package dydx

import (
	"context"
	"math/big"
	"testing"

	"github.com/okex/exchain/libs/dydx/contracts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/stretchr/testify/require"
)

var (
	TopicLogOrderFilled = common.BytesToHash(ethcrypto.Keccak256([]byte("LogOrderFilled(bytes32,bytes32,uint256,Fill)")))
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

	ordersAbi, err := contracts.P1OrdersMetaData.GetAbi()
	require.NoError(t, err)
	//topic, err := abi.MakeTopics([]interface{}{ordersAbi.Events["LogOrderFilled"].ID})
	//require.NoError(t, err)

	var query = ethereum.FilterQuery{
		FromBlock: fromBlockNum,
		ToBlock:   endBlockNum,
		Addresses: []common.Address{
			client.contracts.P1OrdersAddress,
		},
		//Topics: [][]common.Hash{
		//	{TopicLogOrderFilled},
		//},
		Topics: [][]common.Hash{
			{ordersAbi.Events["LogOrderFilled"].ID},
		},
	}
	logs, err := client.ethCli.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	for _, log := range logs {
		l, err := client.contracts.P1Orders.ParseLogOrderFilled(log)
		require.NoError(t, err)
		t.Logf("LogFilled: %+v", l)
	}

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
