package staking

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestHandlerDestroyValidator(t *testing.T) {

	validatorAddr1 := sdk.ValAddress(Addrs[0])
	pk1 := PKs[0]
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitPower)
	keeper := mockKeeper.Keeper
	_ = setInstantUnbondPeriod(keeper, ctx)

	//0. destroy a not exist validator
	destroyValMsg := types.NewMsgDestroyValidator([]byte(validatorAddr1))
	handler := NewHandler(keeper)
	response0 := handler(ctx, destroyValMsg)
	require.False(t, response0.IsOK())
	updates0 := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 0, len(updates0))

	//1. create a validator
	handler = NewHandler(keeper)
	createValMsg := NewTestMsgCreateValidator(validatorAddr1, pk1, DefaultMSD)
	response := handler(ctx, createValMsg)
	require.True(t, response.IsOK())
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	//2. destroy the created validator
	destroyValMsg = types.NewMsgDestroyValidator([]byte(validatorAddr1))
	handler = NewHandler(keeper)
	response2 := handler(ctx, destroyValMsg)
	require.True(t, response2.IsOK())
	updates2 := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates2))
}

type MsgFaked struct {
	Fakeid int
}

func (msg MsgFaked) Route() string { return "token" }

func (msg MsgFaked) Type() string { return "issue" }

// ValidateBasic Implements Msg.
func (msg MsgFaked) ValidateBasic() sdk.Error {
	// check owner
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgFaked) GetSignBytes() []byte {
	return sdk.MustSortJSON([]byte("1"))
}

// GetSigners Implements Msg.
func (msg MsgFaked) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func NewFakeMsg() MsgFaked {
	return MsgFaked{
		Fakeid: 0,
	}
}

func TestHandlerBadMessage(t *testing.T) {

	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitPower)
	keeper := mockKeeper.Keeper
	_ = setInstantUnbondPeriod(keeper, ctx)

	msg := NewFakeMsg()
	handler := NewHandler(keeper)
	r := handler(ctx, msg)
	require.False(t, r.IsOK(), r)
}
