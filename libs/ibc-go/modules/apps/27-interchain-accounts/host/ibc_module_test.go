package host_test

import (
	"fmt"
	"testing"

	types3 "github.com/okex/exchain/libs/tendermint/types"

	types2 "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	banktypes "github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"

	"github.com/gogo/protobuf/proto"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/stretchr/testify/suite"
)

var (
	// TODO: Cosmos-SDK ADR-28: Update crypto.AddressHash() when sdk uses address.Module()
	// https://github.com/cosmos/cosmos-sdk/issues/10225
	//
	// TestAccAddress defines a resuable bech32 address for testing purposes
	addr, _        = sdk.AccAddressFromHex(string(crypto.AddressHash([]byte(icatypes.ModuleName))))
	TestAccAddress = icatypes.GenerateAddress(addr, ibctesting.FirstConnectionID, TestPortID)

	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

	// TestPortID defines a resuable port identifier for testing purposes
	TestPortID, _ = icatypes.NewControllerPortID(TestOwnerAddress)

	// TestVersion defines a resuable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibctesting.FirstConnectionID,
		HostConnectionId:       ibctesting.FirstConnectionID,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
)

type InterchainAccountsTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI
}

func TestICATestSuite(t *testing.T) {
	suite.Run(t, new(InterchainAccountsTestSuite))
}

func (suite *InterchainAccountsTestSuite) SetupTest() {
	types3.UnittestOnlySetMilestoneVenus1Height(-1)
	types3.UnittestOnlySetMilestoneVenus4Height(-1)
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
}

func NewICAPath(chainA, chainB ibctesting.TestChainI) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	return path
}

func RegisterInterchainAccount(endpoint *ibctesting.Endpoint, owner string) error {
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return err
	}

	channelSequence := endpoint.Chain.GetSimApp().GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(endpoint.Chain.GetContext())

	if err := endpoint.Chain.GetSimApp().ICAControllerKeeper.RegisterInterchainAccount(endpoint.Chain.GetContext(), endpoint.ConnectionID, owner, TestVersion); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.Coordinator().CommitBlock(endpoint.Chain)

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID

	return nil
}

// SetupICAPath invokes the InterchainAccounts entrypoint and subsequent channel handshake handlers
func SetupICAPath(path *ibctesting.Path, owner string) error {
	if err := RegisterInterchainAccount(path.EndpointA, owner); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenConfirm(); err != nil {
		return err
	}

	return nil
}

// Test initiating a ChanOpenInit using the host chain instead of the controller chain
// ChainA is the controller chain. ChainB is the host chain
func (suite *InterchainAccountsTestSuite) TestChanOpenInit() {
	suite.SetupTest() // reset
	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	// use chainB (host) for ChanOpenInit
	msg := channeltypes.NewMsgChannelOpenInitV4(path.EndpointB.ChannelConfig.PortID,
		icatypes.Version, channeltypes.ORDERED, []string{path.EndpointB.ConnectionID},
		path.EndpointA.ChannelConfig.PortID, icatypes.ModuleName)
	handler := suite.chainB.GetSimApp().MsgServiceRouter().HandlerWithMsg(msg)
	_, err := handler(suite.chainB.GetContext(), msg)

	suite.Require().Error(err)
}

func (suite *InterchainAccountsTestSuite) TestOnChanOpenTry() {
	var (
		path    *ibctesting.Path
		channel *channeltypes.Channel
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"host submodule disabled", func() {
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), types.NewParams(false, []string{}))
			}, false,
		},
		{
			"success: ICA auth module callback returns error", func() {
				// mock module callback should not be called on host side
				suite.chainB.GetSimApp().ICAAuthModule.IBCApp.OnChanOpenTry = func(ctx sdk.Context, order channeltypes.Order, connectionHops []string,
					portID, channelID string, chanCap *capabilitytypes.Capability,
					counterparty channeltypes.Counterparty, counterpartyVersion string,
				) (string, error) {
					return "", fmt.Errorf("mock ica auth fails")
				}
			}, true,
		},
		{
			"ICA callback fails - invalid channel order", func() {
				channel.Ordering = channeltypes.UNORDERED
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := RegisterInterchainAccount(path.EndpointA, TestOwnerAddress)
			suite.Require().NoError(err)
			path.EndpointB.ChannelID = ibctesting.FirstChannelID

			// default values
			counterparty := channeltypes.NewCounterparty(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
			channel = &channeltypes.Channel{
				State:          channeltypes.TRYOPEN,
				Ordering:       channeltypes.ORDERED,
				Counterparty:   counterparty,
				ConnectionHops: []string{path.EndpointB.ConnectionID},
				Version:        path.EndpointB.ChannelConfig.Version,
			}

			tc.malleate()

			// ensure channel on chainB is set in state
			suite.chainB.GetSimApp().IBCKeeper.V2Keeper.ChannelKeeper.SetChannel(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, *channel)

			module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			chanCap, err := suite.chainB.GetSimApp().GetScopedIBCKeeper().NewCapability(suite.chainB.GetContext(), host.ChannelCapabilityPath(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID))
			suite.Require().NoError(err)

			cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			version, err := cbs.OnChanOpenTry(suite.chainB.GetContext(), channel.Ordering, channel.GetConnectionHops(),
				path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, chanCap, channel.Counterparty, path.EndpointB.ChannelConfig.Version, path.EndpointA.ChannelConfig.Version,
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal("", version)
			}
		})
	}
}

// Test initiating a ChanOpenAck using the host chain instead of the controller chain
// ChainA is the controller chain. ChainB is the host chain
func (suite *InterchainAccountsTestSuite) TestChanOpenAck() {
	suite.SetupTest() // reset
	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	err := RegisterInterchainAccount(path.EndpointA, TestOwnerAddress)
	suite.Require().NoError(err)

	err = path.EndpointB.ChanOpenTry()
	suite.Require().NoError(err)

	// chainA maliciously sets channel to TRYOPEN
	channel := channeltypes.NewChannel(channeltypes.TRYOPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID), []string{path.EndpointA.ConnectionID}, TestVersion)
	suite.chainA.GetSimApp().GetIBCKeeper().ChannelKeeper.SetChannel(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, channel)

	// commit state changes so proof can be created
	//suite.chainA.NextBlock()
	suite.chainA.Coordinator().CommitBlock(suite.chainA)

	path.EndpointB.UpdateClient()

	// query proof from ChainA
	channelKey := host.ChannelKey(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
	proofTry, proofHeight := path.EndpointA.Chain.QueryProof(channelKey)

	// use chainB (host) for ChanOpenAck
	msg := channeltypes.NewMsgChannelOpenAck(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, path.EndpointA.ChannelID, TestVersion, proofTry, proofHeight, icatypes.ModuleName)
	handler := suite.chainB.GetSimApp().MsgServiceRouter().HandlerWithMsg(msg)
	_, err = handler(suite.chainB.GetContext(), msg)

	suite.Require().Error(err)
}

func (suite *InterchainAccountsTestSuite) TestOnChanOpenConfirm() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"host submodule disabled", func() {
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), types.NewParams(false, []string{}))
			}, false,
		},
		{
			"success: ICA auth module callback returns error", func() {
				// mock module callback should not be called on host side
				suite.chainB.GetSimApp().ICAAuthModule.IBCApp.OnChanOpenConfirm = func(
					ctx sdk.Context, portID, channelID string,
				) error {
					return fmt.Errorf("mock ica auth fails")
				}
			}, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			path := NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := RegisterInterchainAccount(path.EndpointA, TestOwnerAddress)
			suite.Require().NoError(err)

			err = path.EndpointB.ChanOpenTry()
			suite.Require().NoError(err)

			err = path.EndpointA.ChanOpenAck()
			suite.Require().NoError(err)

			tc.malleate()

			module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			err = cbs.OnChanOpenConfirm(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// OnChanCloseInit on host (chainB)
func (suite *InterchainAccountsTestSuite) TestOnChanCloseInit() {
	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	err := SetupICAPath(path, TestOwnerAddress)
	suite.Require().NoError(err)

	module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
	suite.Require().NoError(err)

	cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
	suite.Require().True(ok)

	err = cbs.OnChanCloseInit(
		suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID,
	)

	suite.Require().Error(err)
}

func (suite *InterchainAccountsTestSuite) TestOnChanCloseConfirm() {
	var path *ibctesting.Path

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success", func() {}, true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path = NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			tc.malleate() // malleate mutates test data
			module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			err = cbs.OnChanCloseConfirm(
				suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *InterchainAccountsTestSuite) TestOnRecvPacket() {
	var packetData []byte
	testCases := []struct {
		name          string
		malleate      func()
		expAckSuccess bool
	}{
		{
			"success", func() {}, true,
		},
		{
			"host submodule disabled", func() {
				suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), types.NewParams(false, []string{}))
			}, false,
		},
		{
			"success with ICA auth module callback failure", func() {
				suite.chainB.GetSimApp().ICAAuthModule.IBCApp.OnRecvPacket = func(
					ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress,
				) exported.Acknowledgement {
					return channeltypes.NewErrorAcknowledgementV4(fmt.Errorf("failed OnRecvPacket mock callback"))
				}
			}, true,
		},
		{
			"ICA OnRecvPacket fails - cannot unmarshal packet data", func() {
				packetData = []byte("invalid data")
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path := NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)
			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			// send 100stake to interchain account wallet
			amount, _ := sdk.ParseCoinsNormalized("1000wei")
			interchainAccountAddr, _ := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
			bankMsg := &banktypes.MsgSendAdapter{FromAddress: suite.chainB.SenderAccount().GetAddress().String(), ToAddress: interchainAccountAddr, Amount: sdk.CoinsToCoinAdapters(amount)}

			_, err = suite.chainB.SendMsgs(bankMsg)
			suite.Require().NoError(err)

			// build packet data
			msg := &banktypes.MsgSendAdapter{
				FromAddress: interchainAccountAddr,
				ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
				Amount:      sdk.CoinsToCoinAdapters(amount),
			}
			data, err := icatypes.SerializeCosmosTx(suite.chainA.Codec(), []sdk.MsgAdapter{msg})
			suite.Require().NoError(err)

			icaPacketData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: data,
			}
			packetData = icaPacketData.GetBytes()

			// build expected ack
			msgResponseBz, err := proto.Marshal(&banktypes.MsgSendResponseAdapter{})
			suite.Require().NoError(err)

			msgData := &types2.MsgData{
				MsgType: sdk.MsgTypeURL(msg),
				Data:    msgResponseBz,
			}

			expectedTxResponse, err := proto.Marshal(&types2.TxMsgData{
				Data: []*types2.MsgData{msgData},
			})
			suite.Require().NoError(err)

			expectedAck := channeltypes.NewResultAcknowledgement(expectedTxResponse)

			params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
			suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)

			// malleate packetData for test cases
			tc.malleate()

			seq := uint64(1)
			packet := channeltypes.NewPacket(packetData, seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.NewHeight(0, 100), 0)

			tc.malleate()

			module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			ack := cbs.OnRecvPacket(suite.chainB.GetContext(), packet, nil)
			if tc.expAckSuccess {
				suite.Require().True(ack.Success())
				suite.Require().Equal(expectedAck, ack)
			} else {
				suite.Require().False(ack.Success())
			}
		})
	}
}

func (suite *InterchainAccountsTestSuite) TestOnAcknowledgementPacket() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"ICA OnAcknowledgementPacket fails with ErrInvalidChannelFlow", func() {}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path := NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			tc.malleate() // malleate mutates test data

			module, _, err := suite.chainB.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainB.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			cbs, ok := suite.chainB.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			packet := channeltypes.NewPacket(
				[]byte("empty packet data"),
				suite.chainA.SenderAccount().GetSequence(),
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				clienttypes.NewHeight(0, 100),
				0,
			)

			err = cbs.OnAcknowledgementPacket(suite.chainB.GetContext(), packet, []byte("ackBytes"), TestAccAddress)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *InterchainAccountsTestSuite) TestOnTimeoutPacket() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"ICA OnTimeoutPacket fails with ErrInvalidChannelFlow", func() {}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path := NewICAPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			err := SetupICAPath(path, TestOwnerAddress)
			suite.Require().NoError(err)

			tc.malleate() // malleate mutates test data

			module, _, err := suite.chainA.GetSimApp().GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainA.GetContext(), path.EndpointB.ChannelConfig.PortID)
			suite.Require().NoError(err)

			cbs, ok := suite.chainA.GetSimApp().GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			packet := channeltypes.NewPacket(
				[]byte("empty packet data"),
				suite.chainA.SenderAccount().GetSequence(),
				path.EndpointB.ChannelConfig.PortID,
				path.EndpointB.ChannelID,
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				clienttypes.NewHeight(0, 100),
				0,
			)

			err = cbs.OnTimeoutPacket(suite.chainA.GetContext(), packet, TestAccAddress)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *InterchainAccountsTestSuite) fundICAWallet(ctx sdk.Context, portID string, amount sdk.Coins) {
	interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(ctx, ibctesting.FirstConnectionID, portID)
	suite.Require().True(found)

	msgBankSend := &banktypes.MsgSendAdapter{
		FromAddress: suite.chainB.SenderAccount().GetAddress().String(),
		ToAddress:   interchainAccountAddr,
		Amount:      sdk.CoinsToCoinAdapters(amount),
	}

	res, err := suite.chainB.SendMsgs(msgBankSend)
	suite.Require().NotEmpty(res)
	suite.Require().NoError(err)
}

// TestControlAccountAfterChannelClose tests that a controller chain can control a registered interchain account after the currently active channel for that interchain account has been closed
// by opening a new channel on the associated portID
func (suite *InterchainAccountsTestSuite) TestControlAccountAfterChannelClose() {
	// create channel + init interchain account on a particular port
	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)
	err := SetupICAPath(path, TestOwnerAddress)
	suite.Require().NoError(err)

	// check that the account is working as expected
	suite.fundICAWallet(suite.chainB.GetContext(), path.EndpointA.ChannelConfig.PortID, sdk.NewCoins(sdk.NewCoin(sdk.DefaultIbcWei, sdk.NewInt(10000))))
	interchainAccountAddr, found := suite.chainB.GetSimApp().ICAHostKeeper.GetInterchainAccountAddress(suite.chainB.GetContext(), ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID)
	suite.Require().True(found)

	tokenAmt := sdk.NewCoins(sdk.NewCoin(sdk.DefaultIbcWei, sdk.NewInt(5000)))
	msg := &banktypes.MsgSendAdapter{
		FromAddress: interchainAccountAddr,
		ToAddress:   suite.chainB.SenderAccount().GetAddress().String(),
		Amount:      sdk.CoinsToCoinAdapters(tokenAmt),
	}

	data, err := icatypes.SerializeCosmosTx(suite.chainA.GetSimApp().AppCodec(), []sdk.MsgAdapter{msg})
	suite.Require().NoError(err)

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	params := types.NewParams(true, []string{sdk.MsgTypeURL(msg)})
	suite.chainB.GetSimApp().ICAHostKeeper.SetParams(suite.chainB.GetContext(), params)

	chanCap, ok := suite.chainA.GetSimApp().ScopedICAMockKeeper.GetCapability(path.EndpointA.Chain.GetContext(), host.ChannelCapabilityPath(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID))
	suite.Require().True(ok)

	_, err = suite.chainA.GetSimApp().ICAControllerKeeper.SendTx(suite.chainA.GetContext(), chanCap, ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID, icaPacketData, ^uint64(0))
	suite.Require().NoError(err)
	path.EndpointB.UpdateClient()

	// relay the packet
	packetRelay := channeltypes.NewPacket(icaPacketData.GetBytes(), 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.ZeroHeight(), ^uint64(0))
	err = path.RelayPacketV4(packetRelay)
	suite.Require().NoError(err) // relay committed

	// check that the ica balance is updated
	icaAddr, err := sdk.AccAddressFromBech32(interchainAccountAddr)
	suite.Require().NoError(err)

	hasBalance := suite.chainB.GetSimApp().BankKeeper.HasBalance(suite.chainB.GetContext(), icaAddr, sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDecFromInt(sdk.NewInt(5000))})
	suite.Require().True(hasBalance)

	// close the channel
	err = path.EndpointA.SetChannelClosed()
	suite.Require().NoError(err)
	err = path.EndpointB.SetChannelClosed()
	suite.Require().NoError(err)

	// open a new channel on the same port
	path.EndpointA.ChannelID = ""
	path.EndpointB.ChannelID = ""
	suite.coordinator.CreateChannels(path)

	// try to control the interchain account again
	chanCap, ok = suite.chainA.GetSimApp().ScopedICAMockKeeper.GetCapability(path.EndpointA.Chain.GetContext(), host.ChannelCapabilityPath(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID))
	suite.Require().True(ok)

	_, err = suite.chainA.GetSimApp().ICAControllerKeeper.SendTx(suite.chainA.GetContext(), chanCap, ibctesting.FirstConnectionID, path.EndpointA.ChannelConfig.PortID, icaPacketData, ^uint64(0))
	suite.Require().NoError(err)
	path.EndpointB.UpdateClient()

	// relay the packet
	packetRelay = channeltypes.NewPacket(icaPacketData.GetBytes(), 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, clienttypes.ZeroHeight(), ^uint64(0))
	err = path.RelayPacketV4(packetRelay)
	suite.Require().NoError(err) // relay committed

	// check that the ica balance is updated
	hasBalance = suite.chainB.GetSimApp().BankKeeper.HasBalance(suite.chainB.GetContext(), icaAddr, sdk.Coin{Denom: sdk.DefaultIbcWei, Amount: sdk.NewDecFromInt(sdk.NewInt(0))})
	suite.Require().True(hasBalance)
}
