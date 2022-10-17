package types_test

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	banktypes "github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
)

// caseRawBytes defines a helper struct, used for testing codec operations
type caseRawBytes struct {
	name    string
	bz      []byte
	expPass bool
}

var (
	_ sdk.MsgAdapter = mockSdkMsg{}
)

// mockSdkMsg defines a mock struct, used for testing codec error scenarios
type mockSdkMsg struct{}

func (m mockSdkMsg) Route() string {
	return ""
}

func (m mockSdkMsg) Type() string {
	return ""
}

func (m mockSdkMsg) GetSignBytes() []byte {
	return nil
}

// Reset implements sdk.Msg
func (mockSdkMsg) Reset() {
}

// String implements sdk.Msg
func (mockSdkMsg) String() string {
	return ""
}

// ProtoMessage implements sdk.Msg
func (mockSdkMsg) ProtoMessage() {
}

// ValidateBasic implements sdk.Msg
func (mockSdkMsg) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (mockSdkMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (suite *TypesTestSuite) TestSerializeAndDeserializeCosmosTx() {
	testCases := []struct {
		name    string
		msgs    []sdk.MsgAdapter
		expPass bool
	}{
		{
			"single msg",
			[]sdk.MsgAdapter{
				&banktypes.MsgSendAdapter{
					FromAddress: TestOwnerAddress,
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter("bananas", sdk.NewInt(100))},
				},
			},
			true,
		},
		{
			"multiple msgs, same types",
			[]sdk.MsgAdapter{
				&banktypes.MsgSendAdapter{
					FromAddress: TestOwnerAddress,
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter("bananas", sdk.NewInt(100))},
				},
				&banktypes.MsgSendAdapter{
					FromAddress: TestOwnerAddress,
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.CoinAdapters{(sdk.NewCoinAdapter("bananas", sdk.NewInt(200)))},
				},
			},
			true,
		},
		{
			"multiple msgs, different types",
			[]sdk.MsgAdapter{
				&banktypes.MsgSendAdapter{
					FromAddress: TestOwnerAddress,
					ToAddress:   TestOwnerAddress,
					Amount:      sdk.CoinAdapters{(sdk.NewCoinAdapter("bananas", sdk.NewInt(100)))},
				},
			},
			true,
		},
		{
			"unregistered msg type",
			[]sdk.MsgAdapter{
				&mockSdkMsg{},
			},
			false,
		},
		{
			"multiple unregistered msg types",
			[]sdk.MsgAdapter{
				&mockSdkMsg{},
				&mockSdkMsg{},
				&mockSdkMsg{},
			},
			false,
		},
	}
	cdc := suite.chainA.GetSimApp().AppCodec()
	testCasesAny := []caseRawBytes{}

	for _, tc := range testCases {
		bz, err := types.SerializeCosmosTx(cdc, tc.msgs)
		suite.Require().NoError(err, tc.name)

		testCasesAny = append(testCasesAny, caseRawBytes{tc.name, bz, tc.expPass})
	}

	for i, tc := range testCasesAny {
		msgs, err := types.DeserializeCosmosTx(cdc, tc.bz)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].msgs, msgs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}

	// test deserializing unknown bytes
	msgs, err := types.DeserializeCosmosTx(cdc, []byte("invalid"))
	suite.Require().Error(err)
	suite.Require().Empty(msgs)
}
