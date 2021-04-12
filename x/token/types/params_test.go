package types

import (
	"testing"

	"github.com/okex/exchain/x/common"

	"github.com/okex/exchain/x/params"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	common.InitConfig()

	param := DefaultParams()
	expectedString := `Params: 
FeeIssue: 2500.000000000000000000` + common.NativeToken + `
FeeMint: 10.000000000000000000` + common.NativeToken + `
FeeBurn: 10.000000000000000000` + common.NativeToken + `
FeeModify: 0.000000000000000000` + common.NativeToken + `
FeeChown: 10.000000000000000000` + common.NativeToken + `
OwnershipConfirmWindow: 24h0m0s
`

	paramStr := param.String()
	require.EqualValues(t, expectedString, paramStr)

	psp := params.ParamSetPairs{
		{Key: KeyFeeIssue, Value: &param.FeeIssue},
		{Key: KeyFeeMint, Value: &param.FeeMint},
		{Key: KeyFeeBurn, Value: &param.FeeBurn},
		{Key: KeyFeeModify, Value: &param.FeeModify},
		{Key: KeyFeeChown, Value: &param.FeeChown},
		{Key: KeyOwnershipConfirmWindow, Value: &param.OwnershipConfirmWindow},
	}

	for i := range psp {
		require.EqualValues(t, psp[i].Key, param.ParamSetPairs()[i].Key)
		require.EqualValues(t, psp[i].Value, param.ParamSetPairs()[i].Value)
	}

}
