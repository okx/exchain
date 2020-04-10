package types

import (
	"testing"

	"github.com/okex/okchain/x/params"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	param := DefaultParams()
	expectedString := `Params: 
FeeBase: 0.01250000okt
FeeIssue: 20000.00000000okt
FeeMint: 2000.00000000okt
FeeBurn: 10.00000000okt
FeeModify: 0.00000000okt
FeeSend: 0.00000000okt
FeeMultiSend: 0.01000000okt
FeeChown: 10.00000000okt
`
	paramStr := param.String()
	require.EqualValues(t, expectedString, paramStr)

	psp := params.ParamSetPairs{
		{Key: KeyFeeBase, Value: &param.FeeBase},
		{Key: KeyFeeIssue, Value: &param.FeeIssue},
		{Key: KeyFeeMint, Value: &param.FeeMint},
		{Key: KeyFeeBurn, Value: &param.FeeBurn},
		{Key: KeyFeeModify, Value: &param.FeeModify},
		{Key: KeyFeeSend, Value: &param.FeeSend},
		{Key: KeyFeeMultiSend, Value: &param.FeeMultiSend},
		{Key: KeyFeeChown, Value: &param.FeeChown},
	}

	require.EqualValues(t, psp, param.ParamSetPairs())
}
