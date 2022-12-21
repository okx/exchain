package ante_test

import (
	"math/big"
	"testing"

	appante "github.com/okex/exchain/app/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
	"github.com/okex/exchain/libs/ibc-go/testing/mock"
	helpers2 "github.com/okex/exchain/libs/ibc-go/testing/simapp/helpers"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/order"
	"github.com/stretchr/testify/suite"
)

type AnteTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI

	path *ibctesting.Path
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *AnteTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)
	suite.path = ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(suite.path)
}

// TestAnteTestSuite runs all the tests within this package.
func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (suite *AnteTestSuite) TestAnteDecorator() {
	testCases := []struct {
		name     string
		malleate func(suite *AnteTestSuite) []ibcmsg.Msg
		expPass  bool
	}{
		{
			"success on single msg",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				packet := channeltypes.NewPacket([]byte(mock.MockPacketData), 1,
					suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
					suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
					clienttypes.NewHeight(1, 0), 0)

				return []ibcmsg.Msg{channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String())}
			},
			true,
		},
		{
			"success on multiple msgs",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 5; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			true,
		},
		{
			"success on multiple msgs: 1 fresh recv packet",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 5; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					err := suite.path.EndpointA.SendPacket(packet)
					suite.Require().NoError(err)

					// receive all sequences except packet 3
					if i != 3 {
						err = suite.path.EndpointB.RecvPacket(packet)
						suite.Require().NoError(err)
					}

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}

				return msgs
			},
			true,
		},
		{
			"success on multiple mixed msgs",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)
					err := suite.path.EndpointA.SendPacket(packet)
					suite.Require().NoError(err)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)
					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i <= 6; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)
					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					msgs = append(msgs, channeltypes.NewMsgTimeout(packet, uint64(i), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			true,
		},
		{
			"success on multiple mixed msgs: 1 fresh packet of each type",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)
					err := suite.path.EndpointA.SendPacket(packet)
					suite.Require().NoError(err)

					// receive all sequences except packet 3
					if i != 3 {

						err := suite.path.EndpointB.RecvPacket(packet)
						suite.Require().NoError(err)
					}

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)
					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					// receive all acks except ack 2
					if i != 2 {
						err = suite.path.EndpointA.RecvPacket(packet)
						suite.Require().NoError(err)
						err = suite.path.EndpointB.AcknowledgePacket(packet, mock.MockAcknowledgement.Acknowledgement())
						suite.Require().NoError(err)
					}

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i <= 6; i++ {
					height := suite.chainA.LastHeader().GetHeight()
					timeoutHeight := clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+1)
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						timeoutHeight, 0)
					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					// timeout packet
					suite.coordinator.CommitNBlocks(suite.chainA, 3)

					// timeout packets except sequence 5
					if i != 5 {
						suite.path.EndpointB.UpdateClient()
						err = suite.path.EndpointB.TimeoutPacket(packet)
						suite.Require().NoError(err)
					}

					msgs = append(msgs, channeltypes.NewMsgTimeout(packet, uint64(i), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			true,
		},
		{
			"success on multiple mixed msgs: only 1 fresh msg in total",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all packets
					suite.path.EndpointA.SendPacket(packet)
					suite.path.EndpointB.RecvPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all acks
					suite.path.EndpointB.SendPacket(packet)
					suite.path.EndpointA.RecvPacket(packet)
					suite.path.EndpointB.AcknowledgePacket(packet, mock.MockAcknowledgement.Acknowledgement())

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i < 5; i++ {
					height := suite.chainA.LastHeader().GetHeight()
					timeoutHeight := clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+1)
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						timeoutHeight, 0)

					// do not timeout packet, timeout msg is fresh
					suite.path.EndpointB.SendPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgTimeout(packet, uint64(i), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			true,
		},
		{
			"success on single update client msg",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				return []ibcmsg.Msg{&clienttypes.MsgUpdateClient{
					suite.chainB.ChainID(),
					nil,
					suite.chainB.SenderAccount().GetAddress().String(),
				}}
			},
			true,
		},
		{
			"success on multiple update clients",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				return []ibcmsg.Msg{
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()},
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()},
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()}}
			},
			true,
		},
		{
			"success on multiple update clients and fresh packet message",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				msgs := []ibcmsg.Msg{
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()},
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()},
					&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()},
				}

				packet := channeltypes.NewPacket([]byte(mock.MockPacketData), 1,
					suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
					suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
					clienttypes.NewHeight(1, 0), 0)

				return append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
			},
			true,
		},
		{
			"success of tx with different msg type even if all packet messages are redundant",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				msgs := []ibcmsg.Msg{&clienttypes.MsgUpdateClient{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()}}

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all packets
					suite.path.EndpointA.SendPacket(packet)
					suite.path.EndpointB.RecvPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all acks
					suite.path.EndpointB.SendPacket(packet)
					suite.path.EndpointA.RecvPacket(packet)
					suite.path.EndpointB.AcknowledgePacket(packet, mock.MockAcknowledgement.Acknowledgement())

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i < 6; i++ {
					height := suite.chainA.LastHeader().GetHeight()
					timeoutHeight := clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+1)
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						timeoutHeight, 0)

					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					// timeout packet
					suite.coordinator.CommitNBlocks(suite.chainA, 3)

					suite.path.EndpointB.UpdateClient()
					suite.path.EndpointB.TimeoutPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgTimeoutOnClose(packet, uint64(i), []byte("proof"), []byte("channelProof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}

				// append non packet and update message to msgs to ensure multimsg tx should pass
				msgs = append(msgs, &clienttypes.MsgSubmitMisbehaviour{suite.chainB.ChainID(), nil, suite.chainB.SenderAccount().GetAddress().String()})

				return msgs
			},
			true,
		},
		{
			"no success on multiple mixed message: all are redundant",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				var msgs []ibcmsg.Msg

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all packets
					suite.path.EndpointA.SendPacket(packet)
					suite.path.EndpointB.RecvPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all acks
					suite.path.EndpointB.SendPacket(packet)
					suite.path.EndpointA.RecvPacket(packet)
					suite.path.EndpointB.AcknowledgePacket(packet, mock.MockAcknowledgement.Acknowledgement())

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i < 6; i++ {
					height := suite.chainA.LastHeader().GetHeight()
					timeoutHeight := clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+1)
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						timeoutHeight, 0)

					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					// timeout packet
					suite.coordinator.CommitNBlocks(suite.chainA, 3)

					suite.path.EndpointB.UpdateClient()
					suite.path.EndpointB.TimeoutPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgTimeoutOnClose(packet, uint64(i), []byte("proof"), []byte("channelProof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			false,
		},
		{
			"no success if msgs contain update clients and redundant packet messages",
			func(suite *AnteTestSuite) []ibcmsg.Msg {
				msgs := []ibcmsg.Msg{&clienttypes.MsgUpdateClient{Signer: "ex1mnd48anr8jwjjzt347tpdzwfddem08r4d66j3z"}, &clienttypes.MsgUpdateClient{Signer: "ex1mnd48anr8jwjjzt347tpdzwfddem08r4d66j3z"}, &clienttypes.MsgUpdateClient{Signer: "ex1mnd48anr8jwjjzt347tpdzwfddem08r4d66j3z"}}

				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all packets
					suite.path.EndpointA.SendPacket(packet)
					suite.path.EndpointB.RecvPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 1; i <= 3; i++ {
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						clienttypes.NewHeight(1, 0), 0)

					// receive all acks
					suite.path.EndpointB.SendPacket(packet)
					suite.path.EndpointA.RecvPacket(packet)
					suite.path.EndpointB.AcknowledgePacket(packet, mock.MockAcknowledgement.Acknowledgement())

					msgs = append(msgs, channeltypes.NewMsgAcknowledgement(packet, []byte("ack"), []byte("proof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				for i := 4; i < 6; i++ {
					height := suite.chainA.LastHeader().GetHeight()
					timeoutHeight := clienttypes.NewHeight(height.GetRevisionNumber(), height.GetRevisionHeight()+1)
					packet := channeltypes.NewPacket([]byte(mock.MockPacketData), uint64(i),
						suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID,
						suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID,
						timeoutHeight, 0)

					err := suite.path.EndpointB.SendPacket(packet)
					suite.Require().NoError(err)

					// timeout packet
					suite.coordinator.CommitNBlocks(suite.chainA, 3)

					suite.path.EndpointB.UpdateClient()
					suite.path.EndpointB.TimeoutPacket(packet)

					msgs = append(msgs, channeltypes.NewMsgTimeoutOnClose(packet, uint64(i), []byte("proof"), []byte("channelProof"), clienttypes.NewHeight(0, 1), suite.chainB.SenderAccount().GetAddress().String()))
				}
				return msgs
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			// reset suite
			suite.SetupTest()

			k := suite.chainB.App().GetFacadedKeeper()
			//decorator := ante.NewAnteDecorator(k)
			app := suite.chainB.GetSimApp()
			msgs := tc.malleate(suite)

			deliverCtx := suite.chainB.GetContext().WithIsCheckTx(false)
			checkCtx := suite.chainB.GetContext().WithIsCheckTx(true)

			// create multimsg tx
			txBuilder := suite.chainB.TxConfig().NewTxBuilder()
			err := txBuilder.SetMsgs(msgs...)
			ibcTx, err := helpers2.GenTx(
				suite.chainB.TxConfig(),
				msgs,
				sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
				helpers.DefaultGenTxGas,
				suite.chainB.ChainID(),
				[]uint64{suite.chainB.SenderAccount().GetAccountNumber()},
				[]uint64{suite.chainB.SenderAccount().GetSequence()},
				1,
				suite.chainB.SenderAccountPV(),
			)
			antehandler := appante.NewAnteHandler(app.AccountKeeper, app.EvmKeeper, app.SupplyKeeper, validateMsgHook(app.OrderKeeper), app.WasmHandler, k)
			antehandler(deliverCtx, ibcTx, false)
			//_, err = decorator.AnteHandle(deliverCtx, ibcTx, false, next)
			suite.Require().NoError(err, "antedecorator should not error on DeliverTx")
			ibcTx, err = helpers2.GenTx(
				suite.chainB.TxConfig(),
				msgs,
				sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
				helpers.DefaultGenTxGas,
				suite.chainB.ChainID(),
				[]uint64{suite.chainB.SenderAccount().GetAccountNumber()},
				[]uint64{suite.chainB.SenderAccount().GetSequence() + uint64(1)},
				1,
				suite.chainB.SenderAccountPV(),
			)
			//_, err = decorator.AnteHandle(checkCtx, ibcTx, false, next)
			_, err = antehandler(checkCtx, ibcTx, false)
			if tc.expPass {
				suite.Require().NoError(err, "non-strict decorator did not pass as expected")
			} else {
				suite.Require().Error(err, "non-strict antehandler did not return error as expected")
			}
		})
	}
}

func validateMsgHook(orderKeeper order.Keeper) appante.ValidateMsgHandler {
	return func(newCtx sdk.Context, msgs []sdk.Msg) error {

		wrongMsgErr := sdk.ErrUnknownRequest(
			"It is not allowed that a transaction with more than one message contains order or evm message")
		var err error

		for _, msg := range msgs {
			switch assertedMsg := msg.(type) {
			case order.MsgNewOrders:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
				_, err = order.ValidateMsgNewOrders(newCtx, orderKeeper, assertedMsg)
			case order.MsgCancelOrders:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
				err = order.ValidateMsgCancelOrders(newCtx, orderKeeper, assertedMsg)
			case *evmtypes.MsgEthereumTx:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
			}

			if err != nil {
				return err
			}
		}
		return nil
	}
}
