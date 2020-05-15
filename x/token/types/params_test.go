package types

import (
	"testing"

	"github.com/okex/okchain/x/common"

	"github.com/okex/okchain/x/params"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	param := DefaultParams()
	expectedString := `Params: 
FeeIssue: 20000.00000000` + common.NativeToken + `
FeeMint: 2000.00000000` + common.NativeToken + `
FeeBurn: 10.00000000` + common.NativeToken + `
FeeModify: 0.00000000` + common.NativeToken + `
FeeSend: 0.00000000` + common.NativeToken + `
FeeMultiSend: 0.01000000` + common.NativeToken + `
FeeChown: 10.00000000` + common.NativeToken + `
`
	paramStr := param.String()
	require.EqualValues(t, expectedString, paramStr)

	psp := params.ParamSetPairs{
		{Key: KeyFeeIssue, Value: &param.FeeIssue},
		{Key: KeyFeeMint, Value: &param.FeeMint},
		{Key: KeyFeeBurn, Value: &param.FeeBurn},
		{Key: KeyFeeModify, Value: &param.FeeModify},
		{Key: KeyFeeMultiSend, Value: &param.FeeMultiSend},
		{Key: KeyFeeChown, Value: &param.FeeChown},
	}

	require.EqualValues(t, psp, param.ParamSetPairs())
}
