package keeper_test

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	banktypes "github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	icatypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

// TODO,再加
func (suite *KeeperTestSuite) TestSendTx() {
	var (
		path             *ibctesting.Path
		packetData       icatypes.InterchainAccountPacketData
		chanCap          *capabilitytypes.Capability
		timeoutTimestamp uint64
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				interchainAccountAddr, found := suite.chainA.GetSimApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &banktypes.MsgSendAdapter{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
					Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultBondDenom, sdk.NewInt(100))},
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainB.GetSimApp().AppCodec(), []sdk.MsgAdapter{msg})
				suite.Require().NoError(err)

				packetData = icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}
			},
			true,
		},
		{
			"success with multiple sdk.Msg",
			func() {
				interchainAccountAddr, found := suite.chainA.GetSimApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msgsBankSend := []sdk.MsgAdapter{
					&banktypes.MsgSendAdapter{
						FromAddress: interchainAccountAddr,
						ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
						Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultBondDenom, sdk.NewInt(100))},
					},
					&banktypes.MsgSendAdapter{
						FromAddress: interchainAccountAddr,
						ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
						Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultBondDenom, sdk.NewInt(100))},
					},
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainB.GetSimApp().AppCodec(), msgsBankSend)
				suite.Require().NoError(err)

				packetData = icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}
			},
			true,
		},
		{
			"data is nil",
			func() {
				packetData = icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: nil,
				}
			},
			false,
		},
		{
			"active channel not found",
			func() {
				path.EndpointA.ChannelConfig.PortID = "invalid-port-id"
			},
			false,
		},
		{
			"channel does not exist",
			func() {
				suite.chainA.GetSimApp().ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID, "channel-100")
			},
			false,
		},
		{
			"channel in INIT state - optimistic packet sends fail",
			func() {
				channel, found := suite.chainA.GetSimApp().IBCKeeper.V2Keeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				suite.Require().True(found)

				channel.State = channeltypes.INIT
				suite.chainA.GetSimApp().IBCKeeper.V2Keeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)
			},
			false,
		},
		{
			"sendPacket fails - channel closed",
			func() {
				err := path.EndpointA.SetChannelClosed()
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"invalid channel capability provided",
			func() {
				chanCap = nil
			},
			false,
		},
		{
			"timeout timestamp is not in the future",
			func() {
				interchainAccountAddr, found := suite.chainA.GetSimApp().ICAControllerKeeper.GetInterchainAccountAddress(suite.chainA.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
				suite.Require().True(found)

				msg := &banktypes.MsgSendAdapter{
					FromAddress: interchainAccountAddr,
					ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
					Amount:      sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultBondDenom, sdk.NewInt(100))},
				}

				data, err := icatypes.SerializeCosmosTx(suite.chainB.GetSimApp().AppCodec(), []sdk.MsgAdapter{msg})
				suite.Require().NoError(err)

				packetData = icatypes.InterchainAccountPacketData{
					Type: icatypes.EXECUTE_TX,
					Data: data,
				}

				ctx := suite.chainA.GetContext()
				v := &ctx
				timeoutTimestamp = uint64(v.BlockTime().UnixNano())
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.msg, func() {
			suite.SetupTest()             // reset
			timeoutTimestamp = ^uint64(0) // default

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			var ok bool
			chanCap, ok = suite.chainA.GetSimApp().ScopedICAMockKeeper.GetCapability(path.EndpointA.Chain.GetContext(), host.ChannelCapabilityPath(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID))
			suite.Require().True(ok)

			tc.malleate() // malleate mutates test data

			_, err = suite.chainA.GetSimApp().ICAControllerKeeper.SendTx(suite.chainA.GetContext(), chanCap, ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID, packetData, timeoutTimestamp)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	var path *ibctesting.Path

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			tc.malleate() // malleate mutates test data

			packet := channeltypes.NewPacket(
				[]byte{},
				1,
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				clienttypes.NewHeight(0, 100),
				0,
			)

			err = suite.chainA.GetSimApp().ICAControllerKeeper.OnTimeoutPacket(suite.chainA.GetContext(), packet)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
