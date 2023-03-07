package distribution

import (
	"fmt"
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func testMustAccAddressFromBech32(addr string) sdk.AccAddress {
	re, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}
	return re
}

func TestConvertWithdrawDelegatorAllRewardsMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)

	testcases := []struct {
		msgstr  string
		res     types.MsgWithdrawDelegatorAllRewards
		fnCheck func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards)
	}{
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"%s\"}", addr.String()),
			res:    NewMsgWithdrawDelegatorAllRewards(testMustAccAddressFromBech32(addr.String())),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdrawDelegatorAllRewards), res)
			},
		},
		{
			msgstr: `{"delegator_address": "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res:    NewMsgWithdrawDelegatorAllRewards(testMustAccAddressFromBech32("0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9")),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdrawDelegatorAllRewards), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res:    NewMsgWithdrawDelegatorAllRewards(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdrawDelegatorAllRewards), res)
			},
		},
		// error
		{
			msgstr: "123",
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"\"}"),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E\"}"),
			res:    NewMsgWithdrawDelegatorAllRewards(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdrawDelegatorAllRewards) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertWithdrawDelegatorAllRewardsMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}
