package wsclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublicChannel(t *testing.T) {
	//common.SkipSysTestChecker(t)

	okexWS := BaseWS{endPoint: "wss://dexcomreal.bafang.com:8443/"}
	okchainWS := BaseWS{endPoint: "ws://localhost:6666/"}

	okexChannels := []string{
		"dex_spot/candle60s",
		//"dex_spot/ticker",
		//"dex_spot/matches",
		//"dex_spot/optimized_depth",
	}

	okexFilters := []string{
		"tbtc_tusdk",
		//"tbtc_tusdk",
		//"tbtc_tusdk",
		//"tbtc_tusdk",
	}

	okchainChannels := []string{
		"dex_spot/candle60s",
		//"dex_spot/ticker",
		//"dex_spot/matches",
		//"dex_spot/optimized_depth",
	}

	okchainFilters := []string{
		"eos-37c_okt",
		//"eos-37c_okt",
		//"eos-37c_okt",
		//"eos-37c_okt",
	}

	require.True(t, len(okexChannels) == len(okexFilters))
	require.True(t, len(okchainChannels) == len(okchainFilters))
	require.True(t, len(okchainChannels) == len(okexFilters))

	for i := 0; i < len(okchainChannels); i++ {
		CapturePublicChannelNotices(&okexWS, okexChannels[i], okexFilters[i])
		CapturePublicChannelNotices(&okchainWS, okchainChannels[i], okchainFilters[i])

		noDiff := CompareEventsDiff(okexWS.receivedEvents, okchainWS.receivedEvents, true)
		require.True(t, noDiff)
	}
}
