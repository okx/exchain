package channels

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCases = [][]string{
		{
			"BTC",
			"okexchain1_xxx",
			"1",
		},
		{
			"OKT",
			"okexchain1_yyy",
			"2",
		},
		{
			"ETH",
			"okexchain1_zzz",
			"3",
		},
	}
)

func TestPrivateChannel(t *testing.T) {
	for _, testCase := range testCases {
		key := GetSpotAccountKey(testCase[0], testCase[1])
		require.Equal(t, fmt.Sprintf("P3A:dex_spot:account:%s:%s", testCase[0], testCase[1]), key)

		key = GetSpotOrderKey(testCase[0], testCase[1])
		require.Equal(t, fmt.Sprintf("P3A:dex_spot:order:%s:%s", testCase[0], testCase[1]), key)

		key = GetSpotDealKey(testCase[2])
		require.Equal(t, fmt.Sprintf("P3A:dex_spot:deal:%s", testCase[2]), key)
	}
}

func TestPrivateChannel_C(t *testing.T) {
	for _, testCase := range testCases {
		key := GetCSpotAccountKey(testCase[0], testCase[1])
		require.Equal(t, fmt.Sprintf("P3AC:dex_spot:account:%s:%s", testCase[0], testCase[1]), key)

		key = GetCSpotOrderKey(testCase[0], testCase[1])
		require.Equal(t, fmt.Sprintf("P3AC:dex_spot:order:%s:%s", testCase[0], testCase[1]), key)

		key = GetCSpotDealKey(testCase[2])
		require.Equal(t, fmt.Sprintf("P3AC:dex_spot:deal:%s", testCase[2]), key)
	}
}

func TestPublicChannel(t *testing.T) {
	for _, testCase := range testCases {
		key := GetSpotMatchKey(testCase[0])
		require.Equal(t, fmt.Sprintf("P3P:dex_spot:matches:%s", testCase[0]), key)

		key = GetCSpotMatchKey(testCase[0])
		require.Equal(t, fmt.Sprintf("P3C:dex_spot:matches:%s:", testCase[0]), key)
	}
}

func TestDepthChannel(t *testing.T) {
	testCases := []string{
		"BTC-OKB",
		"BTC-USD",
	}
	for _, testCase := range testCases {
		key := GetSpotDepthKey(testCase)
		require.Equal(t, fmt.Sprintf("P3D:dex_spot:depth:%s", testCase), key)
	}
	for _, testCase := range testCases {
		key := GetCSpotDepthKey(testCase)
		require.Equal(t, fmt.Sprintf("P3DC:dex_spot:depth:%s:", testCase), key)
	}
}

func TestGetKeyNoArgs(t *testing.T) {
	testCase := []string{
		"testchan",
		"testsrv",
		"testop",
	}
	expected := fmt.Sprintf("%s:%s:%s", testCase[0], testCase[1], testCase[2])
	got := getKey(testCase[0], testCase[1], testCase[2], []string{})
	require.Equal(t, expected, got)

	got = getKey(testCase[0], testCase[1], testCase[2], nil)
	require.Equal(t, expected, got)
}

func TestGetSpotKey(t *testing.T) {
	tChan, tSrv, tOp := "testchan", "dex_spot", "testop"
	spotKey := GetSpotKey(tChan, tOp, nil)

	require.Equal(t, fmt.Sprintf("%s:%s:%s", tChan, tSrv, tOp), spotKey)
}

func TestSpotMetaKey(t *testing.T) {
	require.Equal(t, "P3K:dex_spot:instruments", GetSpotMetaKey())
	require.Equal(t, "P3KC:dex_spot:instruments", GetCSpotMetaKey())
}
