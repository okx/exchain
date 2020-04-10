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
	valOpAddrs, valConsPks, valConsAddrs := keeper.GetTestAddrs()
	ctx, ak, _, k, sk, _, supplyKeeper := keeper.CreateTestInputAdvanced(t, false, 1000)
	dh := NewHandler(k)

	// create one validator
	sh := staking.NewHandler(sk)
	skMsg := staking.NewMsgCreateValidator(valOpAddrs[0], valConsPks[0],
		staking.Description{}, keeper.NewTestDecCoin(1, 0))
	require.True(t, sh(ctx, skMsg).IsOK())

	//send 1okt fee
	feeCollector := supplyKeeper.GetModuleAccount(ctx, k.GetFeeCollectorName())
	err := feeCollector.SetCoins(keeper.NewTestDecCoins(1, 0))
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
	// crate votes info and allocate tokens
	abciVal := abci.Validator{Address: valConsPks[0].Address(), Power: 1}
	votes := []abci.VoteInfo{{Validator: abciVal, SignedLastBlock: true}}
	k.AllocateTokens(ctx, 100, valConsAddrs[0], votes)

	//send withdraw-commission msgWithdrawValCommission
	msgWithdrawValCommission := types.NewMsgWithdrawValidatorCommission(valOpAddrs[0])
	require.True(t, dh(ctx, msgWithdrawValCommission).IsOK())
	require.False(t, dh(ctx, msgWithdrawValCommission).IsOK())

	//send set-withdraw-address msgSetWithdrawAddress
	msgSetWithdrawAddress := types.NewMsgSetWithdrawAddress(keeper.TestAddrs[0], keeper.TestAddrs[1])
	require.True(t, dh(ctx, msgSetWithdrawAddress).IsOK())
	k.SetWithdrawAddrEnabled(ctx, false)
	require.False(t, dh(ctx, msgSetWithdrawAddress).IsOK())
	msgSetWithdrawAddress = types.NewMsgSetWithdrawAddress(keeper.TestAddrs[0],
		supplyKeeper.GetModuleAddress(ModuleName))
	require.False(t, dh(ctx, msgSetWithdrawAddress).IsOK())

	//send unknown msgWithdrawValCommission
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
