package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func testCode(t *testing.T, err sdk.Error, expectedCode uint32) {
	if expectedCode != 0 {
		require.NotNil(t, err)
	}else {
		require.Nil(t, err)
	}
}

func TestMsgCreatePool(t *testing.T) {
	negMinLockAmount := sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec())
	negMinLockAmount.Amount = sdk.NewDec(-1)
	tests := []struct {
		owner         sdk.AccAddress
		poolName      string
		minLockAmount sdk.SysCoin
		yieldedSymbol string
		errCode       uint32
	}{
		{sdk.AccAddress{0x1}, "pool", sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec()), "wwb", sdk.CodeOK},
		{nil, "pool", sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec()), "wwb", sdk.CodeInvalidAddress},
		{sdk.AccAddress{0x1}, "", sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec()), "wwb", CodeInvalidInput},
		{sdk.AccAddress{0x1}, "pool", negMinLockAmount, "wwb", sdk.CodeOK},
		{sdk.AccAddress{0x1}, "pool", sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec()), "", CodeInvalidInput},
	}

	for _, test := range tests {
		msg := NewMsgCreatePool(test.owner, test.poolName, test.minLockAmount, test.yieldedSymbol)
		require.Equal(t, createPoolMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.owner}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}

func TestMsgDestroyPool(t *testing.T) {
	tests := []struct {
		owner    sdk.AccAddress
		poolName string
		errCode  uint32
	}{
		{sdk.AccAddress{0x1}, "pool", sdk.CodeOK},
		{nil, "pool", CodeInvalidAddress},
		{sdk.AccAddress{0x1}, "", CodeInvalidInput},
	}

	for _, test := range tests {
		msg := NewMsgDestroyPool(test.owner, test.poolName)
		require.Equal(t, destroyPoolMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.owner}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}

func TestMsgProvide(t *testing.T) {
	tests := []struct {
		poolName         string
		owner            sdk.AccAddress
		amount           sdk.SysCoin
		yieldPerBlock    sdk.Dec
		startBlockHeight int64
		errCode          uint32
	}{
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(10),
			1,
			sdk.CodeOK,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(100),
			1,
			sdk.CodeOK,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(1000),
			1,
			CodeInvalidInput,
		},
		{
			"pool",
			nil,
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(10),
			1,
			CodeInvalidAddress,
		},
		{
			"",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(10),
			1,
			CodeInvalidInput,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(0)),
			sdk.NewDec(10),
			1,
			CodeInvalidInput,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(0),
			1,
			CodeInvalidInput,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.NewDec(10),
			0,
			CodeInvalidInput,
		},
	}

	for _, test := range tests {
		msg := NewMsgProvide(test.poolName, test.owner, test.amount, test.yieldPerBlock, test.startBlockHeight)
		require.Equal(t, provideMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.owner}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}

func TestMsgLock(t *testing.T) {
	tests := []struct {
		poolName string
		addr     sdk.AccAddress
		amount   sdk.SysCoin
		errCode  uint32
	}{
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.CodeOK,
		},
		{
			"pool",
			nil,
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidAddress,
		},
		{
			"",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidInput,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(0)),
			CodeInvalidInput,
		},
	}

	for _, test := range tests {
		msg := NewMsgLock(test.poolName, test.addr, test.amount)
		require.Equal(t, lockMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.addr}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}

func TestMsgUnlock(t *testing.T) {
	tests := []struct {
		poolName string
		addr     sdk.AccAddress
		amount   sdk.SysCoin
		errCode  uint32
	}{
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.CodeOK,
		},
		{
			"pool",
			nil,
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidAddress,
		},
		{
			"",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidInput,
		},
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(0)),
			CodeInvalidInput,
		},
	}

	for _, test := range tests {
		msg := NewMsgUnlock(test.poolName, test.addr, test.amount)
		require.Equal(t, unlockMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.addr}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}

func TestMsgClaim(t *testing.T) {
	tests := []struct {
		poolName string
		addr     sdk.AccAddress
		amount   sdk.SysCoin
		errCode  uint32
	}{
		{
			"pool",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			sdk.CodeOK,
		},
		{
			"pool",
			nil,
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidAddress,
		},
		{
			"",
			sdk.AccAddress{0x1},
			sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			CodeInvalidInput,
		},
	}

	for _, test := range tests {
		msg := NewMsgClaim(test.poolName, test.addr)
		require.Equal(t, claimMsgType, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, []sdk.AccAddress{test.addr}, msg.GetSigners())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		err := msg.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}
	}
}
