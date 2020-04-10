package distribution

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/keeper"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestHandler(t *testing.T) {
	valOpAddrs, valConsPks, valConsAddrs := keeper.GetAddrs()
	ctx, ak, k, sk, supplyKeeper := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)
	dh := NewHandler(k)

	// create one validator
	skMsg := staking.NewMsgCreateValidator(valOpAddrs[0], valConsPks[0], staking.Description{}, keeper.NewDecCoin(1))
	require.True(t, sh(ctx, skMsg).IsOK())

	//send 1okt fee
	feeCollector := supplyKeeper.GetModuleAccount(ctx, k.GetFeeCollectorName())
	require.NotNil(t, feeCollector)
	err := feeCollector.SetCoins(keeper.NewDecCoins(1, 0))
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
	// crate votes info and allocate tokens
	abciVal := abci.Validator{Address: valConsPks[0].Address(), Power: 1}
	votes := []abci.VoteInfo{{Validator: abciVal, SignedLastBlock: true}}
	k.AllocateTokens(ctx, 100, valConsAddrs[0], votes)

	//send withdraw-comssion msg
	msg := types.NewMsgWithdrawValidatorCommission(valOpAddrs[0])
	require.True(t, dh(ctx, msg).IsOK())
	require.False(t, dh(ctx, msg).IsOK())

	//send set-withdraw-address msg
	msg1 := types.NewMsgSetWithdrawAddress(keeper.DelAddr1, keeper.DelAddr2)
	require.True(t, dh(ctx, msg1).IsOK())
	msg1 = types.NewMsgSetWithdrawAddress(keeper.DelAddr1, supplyKeeper.GetModuleAddress(ModuleName))
	require.False(t, dh(ctx, msg1).IsOK())
	k.SetWithdrawAddrEnabled(ctx, false)
	msg1 = types.NewMsgSetWithdrawAddress(keeper.DelAddr1, keeper.DelAddr2)
	require.False(t, dh(ctx, msg1).IsOK())

	//send unknown msg
	fakeMsg := NewMsgFake()
	require.False(t, dh(ctx, fakeMsg).IsOK())
}

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgFake struct {
}

func NewMsgFake() MsgFake {
	return MsgFake{}
}

func (msg MsgFake) Route() string { return "" }
func (msg MsgFake) Type() string  { return "" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgFake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

// get the bytes for the message signer to sign on
func (msg MsgFake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgFake) ValidateBasic() sdk.Error {
	return nil
}
