package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestValidatorTestEquivalent(t *testing.T) {
	val1 := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)
	val2 := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)

	ok := val1.TestEquivalent(val2)
	require.True(t, ok)

	val2 = NewValidator(valAddr2, pk2, Description{}, DefaultMinSelfDelegation)

	ok = val1.TestEquivalent(val2)
	require.False(t, ok)

	// MarshalYAML
	data, err := val1.MarshalYAML()
	require.Nil(t, err)
	require.Contains(t, data, "operatoraddress", data)

	data, err = val2.Standardize().MarshalYAML()
	require.Nil(t, err)
	require.Contains(t, data, "Operator Address")

}

func TestValidators(t *testing.T) {
	val1 := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)
	val2 := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)
	valdators := Validators{val1, val2}

	vaStr := valdators.String()
	require.True(t, len(vaStr) > 0, vaStr)
	stdVas := valdators.Standardize()
	require.NotNil(t, stdVas)
	iVas := valdators.ToSDKValidators()
	require.NotNil(t, iVas)
	require.True(t, len(iVas) == len(valdators))

}

func getFixSizeString(size int) string {
	bytes := make([]byte, size)
	return string(bytes)
}

func TestUpdateDescription(t *testing.T) {
	d1 := Description{Website: "https://validator.cosmos", Details: "Test validator"}
	d2 := Description{
		Moniker:  DoNotModifyDesc,
		Identity: DoNotModifyDesc,
		Website:  DoNotModifyDesc,
		Details:  DoNotModifyDesc,
	}
	d3 := Description{Moniker: "", Identity: "", Website: "", Details: ""}
	d4 := Description{Moniker: getFixSizeString(MaxMonikerLength + 1), Identity: "", Website: "", Details: ""}
	d5 := Description{Moniker: "", Identity: getFixSizeString(MaxIdentityLength + 1), Website: "", Details: ""}
	d6 := Description{Moniker: "", Identity: "", Website: getFixSizeString(MaxWebsiteLength + 1), Details: ""}
	d7 := Description{Moniker: "", Identity: "", Website: "", Details: getFixSizeString(MaxDetailsLength + 1)}

	tests := []struct {
		name       string
		fromDesc   Description
		toDesc     Description
		expectPass bool
	}{
		{"success update1", d1, d2, true},
		{"success update2", d1, d3, true},
		{"fail update of MaxMonikerLength", d1, d4, false},
		{"fail update of MaxIdentityLength", d1, d5, false},
		{"fail update of MaxWebsiteLength", d1, d6, false},
		{"fail update of MaxDetailsLength", d1, d7, false},
	}

	for _, tc := range tests {
		_, err := tc.fromDesc.UpdateDescription(tc.toDesc)

		if tc.expectPass {
			require.Nil(t, err, "test: %v", tc.name)
		} else {
			require.NotNil(t, err, "test: %v", tc.name)
		}
	}
}

func TestABCIValidatorUpdate(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)

	abciVal := validator.ABCIValidatorUpdate()
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.ConsPubKey), abciVal.PubKey)
	require.Equal(t, validator.BondedTokens().Int64(), abciVal.Power)
}

func TestABCIValidatorUpdateZero(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)

	abciVal := validator.ABCIValidatorUpdateZero()
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.ConsPubKey), abciVal.PubKey)
	require.Equal(t, int64(0), abciVal.Power)
}

func TestShareTokens(t *testing.T) {
	validator := Validator{
		OperatorAddress: valAddr1,
		ConsPubKey:      pk1,
		Status:          sdk.Bonded,
		Tokens:          sdk.NewInt(100),
		DelegatorShares: sdk.NewDec(100),
	}
	assert.True(sdk.DecEq(t, sdk.NewDec(50), validator.TokensFromShares(sdk.NewDec(50))))

	validator.Tokens = sdk.NewInt(50)
	assert.True(sdk.DecEq(t, sdk.NewDec(25), validator.TokensFromShares(sdk.NewDec(50))))
	assert.True(sdk.DecEq(t, sdk.NewDec(5), validator.TokensFromShares(sdk.NewDec(10))))
}

func TestValidatorMarshalUnmarshalJSON(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)
	js, err := codec.Cdc.MarshalJSON(validator)
	require.NoError(t, err)
	require.NotEmpty(t, js)
	require.Contains(t, string(js), "\"consensus_pubkey\":\"okexchainvalconspub")
	got := &Validator{}
	err = codec.Cdc.UnmarshalJSON(js, got)
	assert.NoError(t, err)
	assert.Equal(t, validator, *got)
}

func TestValidatorSetInitialCommission(t *testing.T) {
	val := NewValidator(valAddr1, pk1, Description{}, DefaultMinSelfDelegation)
	testCases := []struct {
		validator   Validator
		commission  Commission
		expectedErr bool
	}{
		{val, NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()), false},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDec(15000000000), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec(), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.NewDecWithPrec(2, 1), sdk.NewDecWithPrec(1, 1), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1)), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(2, 1)), true},
	}

	for i, tc := range testCases {
		val, err := tc.validator.SetInitialCommission(tc.commission)

		if tc.expectedErr {
			require.Error(t, err,
				"expected error for test case #%d with commission: %s", i, tc.commission,
			)
		} else {
			require.NoError(t, err,
				"unexpected error for test case #%d with commission: %s", i, tc.commission,
			)
			require.Equal(t, tc.commission, val.Commission,
				"invalid validator commission for test case #%d with commission: %s", i, tc.commission,
			)
		}
	}
}
