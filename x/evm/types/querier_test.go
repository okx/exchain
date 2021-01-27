package types

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestQueryString(t *testing.T) {
	const (
		balance                        = "1024.1024okt"
		number                         = int64(1024)
		nonce                          = uint64(1024)
		expectedQueryResBalanceStr     = "1024.1024okt"
		expectedQueryResBlockNumberStr = "1024"
		expectedBytesStr               = "test bytes"
		expectedQueryResNonceStr       = "1024"
	)

	bytes := []byte("test bytes")
	queryResBalance := QueryResBalance{balance}
	require.True(t, strings.EqualFold(expectedQueryResBalanceStr, queryResBalance.String()))

	queryResBlockNumber := QueryResBlockNumber{number}
	require.True(t, strings.EqualFold(expectedQueryResBlockNumberStr, queryResBlockNumber.String()))

	queryResCode := QueryResCode{bytes}
	require.True(t, strings.EqualFold(expectedBytesStr, queryResCode.String()))

	queryResStorage := QueryResStorage{bytes}
	require.True(t, strings.EqualFold(expectedBytesStr, queryResStorage.String()))

	queryResNonce := QueryResNonce{nonce}
	require.True(t, strings.EqualFold(expectedQueryResNonceStr, queryResNonce.String()))

	bloom := ethtypes.BytesToBloom(bytes)
	queryBloomFilter := QueryBloomFilter{bloom}
	require.True(t, strings.EqualFold(string(bloom[:]), queryBloomFilter.String()))
}

func TestQueryETHLogs_String(t *testing.T) {
	const expectedQueryETHLogsStr = `{0x0000000000000000000000000000000000000000 [] [1 2 3 4] 9 0x0000000000000000000000000000000000000000000000000000000000000000 0 0x0000000000000000000000000000000000000000000000000000000000000000 0 false}
{0x0000000000000000000000000000000000000000 [] [5 6 7 8] 10 0x0000000000000000000000000000000000000000000000000000000000000000 0 0x0000000000000000000000000000000000000000000000000000000000000000 0 false}
`
	logs := []*ethtypes.Log{
		{
			Data:        []byte{1, 2, 3, 4},
			BlockNumber: 9,
		},
		{
			Data:        []byte{5, 6, 7, 8},
			BlockNumber: 10,
		},
	}

	require.True(t, strings.EqualFold(expectedQueryETHLogsStr, QueryETHLogs{logs}.String()))
}
