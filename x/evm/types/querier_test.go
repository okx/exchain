package types

import (
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
		expectedQueryBloomFilterStr    = "test bytes                                                                                                                                                                                                                                                      "
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

	//queryBloomFilter := QueryBloomFilter{ethtypes.BytesToBloom(bytes)}
	//require.True(t, strings.EqualFold(expectedQueryBloomFilterStr, queryBloomFilter.String()))

}
