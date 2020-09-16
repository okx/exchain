package types

import (
	"testing"

	"github.com/okex/okexchain/x/common"

	"github.com/okex/okexchain/x/params"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	param := DefaultParams()
	expectedString := `Params: 
FeeIssue: 2500.00000000` + common.NativeToken + `
FeeMint: 10.00000000` + common.NativeToken + `
FeeBurn: 10.00000000` + common.NativeToken + `
FeeModify: 0.00000000` + common.NativeToken + `
FeeChown: 10.00000000` + common.NativeToken + `
ConfirmPeriod: 24h0m0s
`

	paramStr := param.String()
	require.EqualValues(t, expectedString, paramStr)

	psp := params.ParamSetPairs{
		{Key: KeyFeeIssue, Value: &param.FeeIssue},
		{Key: KeyFeeMint, Value: &param.FeeMint},
		{Key: KeyFeeBurn, Value: &param.FeeBurn},
		{Key: KeyFeeModify, Value: &param.FeeModify},
		{Key: KeyFeeChown, Value: &param.FeeChown},
		{Key: KeyOwnershipConfirmWindow, Value: &param.ConfirmPeriod},
	}

	require.EqualValues(t, psp, param.ParamSetPairs())
}
