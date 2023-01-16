package keeper_test

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	abcitypes "github.com/okex/exchain/libs/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestIncentivizePacketEvent() {
	var (
		expRecvFees    sdk.CoinAdapters
		expAckFees     sdk.CoinAdapters
		expTimeoutFees sdk.CoinAdapters
	)

	suite.coordinator.Setup(suite.path)

	fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)
	msg := types.NewMsgPayPacketFee(
		fee,
		suite.path.EndpointA.ChannelConfig.PortID,
		suite.path.EndpointA.ChannelID,
		suite.chainA.SenderAccount().GetAddress().String(),
		nil,
	)

	expRecvFees = expRecvFees.Add(fee.RecvFee...)
	expAckFees = expAckFees.Add(fee.AckFee...)
	expTimeoutFees = expTimeoutFees.Add(fee.TimeoutFee...)

	result, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err)

	var incentivizedPacketEvent abcitypes.Event
	for _, event := range result.Events {
		if event.Type == types.EventTypeIncentivizedPacket {
			incentivizedPacketEvent = abcitypes.Event(event)
		}
	}

	for _, attr := range incentivizedPacketEvent.Attributes {
		switch string(attr.Key) {
		case types.AttributeKeyRecvFee:
			suite.Require().Equal(expRecvFees.String(), string(attr.Value))

		case types.AttributeKeyAckFee:
			suite.Require().Equal(expAckFees.String(), string(attr.Value))

		case types.AttributeKeyTimeoutFee:
			suite.Require().Equal(expTimeoutFees.String(), string(attr.Value))
		}
	}

	// send the same messages again a few times
	for i := 0; i < 3; i++ {
		expRecvFees = expRecvFees.Add(fee.RecvFee...)
		expAckFees = expAckFees.Add(fee.AckFee...)
		expTimeoutFees = expTimeoutFees.Add(fee.TimeoutFee...)

		result, err = suite.chainA.SendMsgs(msg)
		suite.Require().NoError(err)
	}

	for _, event := range result.Events {
		if event.Type == types.EventTypeIncentivizedPacket {
			incentivizedPacketEvent = abcitypes.Event(event)
		}
	}

	for _, attr := range incentivizedPacketEvent.Attributes {
		switch string(attr.Key) {
		case types.AttributeKeyRecvFee:
			suite.Require().Equal(expRecvFees.String(), string(attr.Value))

		case types.AttributeKeyAckFee:
			suite.Require().Equal(expAckFees.String(), string(attr.Value))

		case types.AttributeKeyTimeoutFee:
			suite.Require().Equal(expTimeoutFees.String(), string(attr.Value))
		}
	}
}
