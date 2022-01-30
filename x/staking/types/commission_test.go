package types

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/tendermint/go-amino"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func TestCommissionValidate(t *testing.T) {
	testCases := []struct {
		input     Commission
		expectErr bool
	}{
		// invalid commission; max rate < 0%
		{NewCommission(sdk.ZeroDec(), sdk.MustNewDecFromStr("-1.00"), sdk.ZeroDec()), true},
		// invalid commission; max rate > 100%
		{NewCommission(sdk.ZeroDec(), sdk.MustNewDecFromStr("2.00"), sdk.ZeroDec()), true},
		// invalid commission; rate < 0%
		{NewCommission(sdk.MustNewDecFromStr("-1.00"), sdk.ZeroDec(), sdk.ZeroDec()), true},
		// invalid commission; rate > max rate
		{NewCommission(sdk.MustNewDecFromStr("0.75"), sdk.MustNewDecFromStr("0.50"), sdk.ZeroDec()), true},
		// invalid commission; max change rate < 0%
		{NewCommission(sdk.OneDec(), sdk.OneDec(), sdk.MustNewDecFromStr("-1.00")), true},
		// invalid commission; max change rate > max rate
		{NewCommission(sdk.OneDec(), sdk.MustNewDecFromStr("0.75"), sdk.MustNewDecFromStr("0.90")), true},
		// valid commission
		{NewCommission(sdk.MustNewDecFromStr("0.20"), sdk.OneDec(), sdk.MustNewDecFromStr("0.10")), false},
	}

	for i, tc := range testCases {
		err := tc.input.Validate()
		require.Equal(t, tc.expectErr, err != nil, "unexpected result; tc #%d, input: %v", i, tc.input)
	}
}

func TestCommissionValidateNewRate(t *testing.T) {
	now := time.Now().UTC()
	c1 := NewCommission(sdk.MustNewDecFromStr("0.40"), sdk.MustNewDecFromStr("0.80"), sdk.MustNewDecFromStr("0.10"))
	c1.UpdateTime = now

	testCases := []struct {
		input     Commission
		newRate   sdk.Dec
		blockTime time.Time
		expectErr bool
	}{
		// invalid new commission rate; last update < 24h ago
		{c1, sdk.MustNewDecFromStr("0.50"), now, true},
		// invalid new commission rate; new rate < 0%
		{c1, sdk.MustNewDecFromStr("-1.00"), now.Add(48 * time.Hour), true},
		// invalid new commission rate; new rate > max rate
		{c1, sdk.MustNewDecFromStr("0.90"), now.Add(48 * time.Hour), true},
		// invalid new commission rate; new rate > max change rate
		{c1, sdk.MustNewDecFromStr("0.60"), now.Add(48 * time.Hour), true},
		// valid commission
		{c1, sdk.MustNewDecFromStr("0.50"), now.Add(48 * time.Hour), false},
		// valid commission
		{c1, sdk.MustNewDecFromStr("0.10"), now.Add(48 * time.Hour), false},
	}

	for i, tc := range testCases {
		err := tc.input.ValidateNewRate(tc.newRate, tc.blockTime)
		require.Equal(
			t, tc.expectErr, err != nil,
			"unexpected result; tc #%d, input: %v, newRate: %s, blockTime: %s",
			i, tc.input, tc.newRate, tc.blockTime,
		)
	}
}

func TestCommissionAmino(t *testing.T) {
	testCases := []Commission{
		{},
		{
			CommissionRates{sdk.NewDec(1), sdk.NewDec(2), sdk.NewDec(3)}, time.Now(),
		},
	}
	cdc := amino.NewCodec()
	for _, commission := range testCases {
		bz, err := cdc.MarshalBinaryBare(commission)
		require.NoError(t, err)

		var newCommission Commission
		err = cdc.UnmarshalBinaryBare(bz, &newCommission)
		require.NoError(t, err)

		var newCommission2 Commission
		err = newCommission2.UnmarshalFromAmino(cdc, bz)
		require.NoError(t, err)

		require.Equal(t, newCommission, newCommission2)
	}
}

func TestCommissionRatesAmino(t *testing.T) {
	testCases := []CommissionRates{
		{},
		{
			sdk.Dec{new(big.Int)},
			sdk.Dec{new(big.Int)},
			sdk.Dec{new(big.Int)},
		},
		{
			sdk.Dec{big.NewInt(1)},
			sdk.Dec{big.NewInt(10)},
			sdk.Dec{big.NewInt(100)},
		},
		{
			sdk.Dec{big.NewInt(math.MinInt64)},
			sdk.Dec{big.NewInt(math.MinInt64)},
			sdk.Dec{big.NewInt(math.MinInt64)},
		},
		{
			sdk.Dec{big.NewInt(math.MaxInt64)},
			sdk.Dec{big.NewInt(math.MaxInt64)},
			sdk.Dec{big.NewInt(math.MaxInt64)},
		},
		{
			sdk.Dec{big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))},
			sdk.Dec{big.NewInt(0).Add(big.NewInt(math.MaxInt64), big.NewInt(math.MaxInt64))},
			sdk.Dec{big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(2))},
		},
	}
	cdc := amino.NewCodec()
	for _, commission := range testCases {
		bz, err := cdc.MarshalBinaryBare(commission)
		require.NoError(t, err)

		var newCommission CommissionRates
		err = cdc.UnmarshalBinaryBare(bz, &newCommission)
		require.NoError(t, err)

		var newCommission2 CommissionRates
		err = newCommission2.UnmarshalFromAmino(cdc, bz)
		require.NoError(t, err)

		require.Equal(t, newCommission, newCommission2)
	}
}
