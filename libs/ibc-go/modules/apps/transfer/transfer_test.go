package transfer_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI
	chainC ibctesting.TestChainI
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewTransferPath(chainA, chainB ibctesting.TestChainI) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
func getBalance(coins sdk.Coins, denom string) *sdk.Coin {
	for _, coin := range coins {
		if coin.Denom == denom {
			return &coin
		}
	}
	return nil
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	pathA2B := NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(pathA2B)

	//	originalBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount().GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(0, 110)

	amount, ok := sdk.NewIntFromString("92233720368547758080") // 2^63 (one above int64)
	suite.Require().True(ok)
	transferAmountDec := sdk.NewDecFromIntWithPrec(amount, 0)
	coinToSendToB := sdk.NewCoin(sdk.DefaultIbcWei, amount)

	// send from chainA to chainB
	msg := types.NewMsgTransfer(pathA2B.EndpointA.ChannelConfig.PortID, pathA2B.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount().GetAddress(), suite.chainB.SenderAccount().GetAddress().String(), timeoutHeight, 0)

	_, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	// relay send

	fungibleTokenPacket := types.NewFungibleTokenPacketData(coinToSendToB.Denom, transferAmountDec.BigInt().String(), suite.chainA.SenderAccount().GetAddress().String(), suite.chainB.SenderAccount().GetAddress().String())
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, pathA2B.EndpointA.ChannelConfig.PortID, pathA2B.EndpointA.ChannelID, pathA2B.EndpointB.ChannelConfig.PortID, pathA2B.EndpointB.ChannelID, timeoutHeight, 0)

	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = pathA2B.RelayPacket(packet, ack.Acknowledgement())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	voucherDenomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), sdk.DefaultIbcWei))
	//balance := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount().GetAddress(), voucherDenomTrace.IBCDenom())

	balanceB := suite.chainB.GetSimApp().BankKeeper.GetCoins(suite.chainB.GetContext(), suite.chainB.SenderAccount().GetAddress())

	denomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(pathA2B.EndpointB.ChannelConfig.PortID, pathA2B.EndpointB.ChannelID, sdk.DefaultIbcWei))
	coinSentFromAToB := sdk.NewCoin(denomTrace.IBCDenom(), amount)
	suite.Require().Equal(coinSentFromAToB, balanceB[0])

	// setup between chainB to chainC
	// NOTE:
	// pathBtoC.EndpointA = endpoint on chainB
	// pathBtoC.EndpointB = endpoint on chainC
	pathBtoC := NewTransferPath(suite.chainB, suite.chainC)
	suite.coordinator.Setup(pathBtoC)

	// send from chainB to chainC
	msg = types.NewMsgTransfer(pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID, coinSentFromAToB, suite.chainB.SenderAccount().GetAddress(), suite.chainC.SenderAccount().GetAddress().String(), timeoutHeight, 0)

	_, err = suite.chainB.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPathB2C := types.GetPrefixedDenom(pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, voucherDenomTrace.GetFullDenomPath())
	denom := types.ParseDenomTrace(fullDenomPathB2C).IBCDenom()
	fungibleTokenPacket = types.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), transferAmountDec.BigInt().String(), suite.chainB.SenderAccount().GetAddress().String(), suite.chainC.SenderAccount().GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID, pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, timeoutHeight, 0)
	err = pathBtoC.RelayPacket(packet, ack.Acknowledgement())
	suite.Require().NoError(err) // relay committed

	coinSentFromBToC := sdk.NewCoin(denom, amount)
	balanceB = suite.chainB.GetSimApp().BankKeeper.GetCoins(suite.chainB.GetContext(), suite.chainB.SenderAccount().GetAddress())
	balanceC := suite.chainC.GetSimApp().BankKeeper.GetCoins(suite.chainC.GetContext(), suite.chainC.SenderAccount().GetAddress())

	// check that the balance is updated on chainC
	suite.Require().Equal(coinSentFromBToC, *getBalance(balanceC, denom))

	// check that balance on chain B is empty
	suite.Require().Nil(getBalance(balanceB, denom))

	// send from chainC back to chainB
	msg = types.NewMsgTransfer(pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, coinSentFromBToC, suite.chainC.SenderAccount().GetAddress(), suite.chainB.SenderAccount().GetAddress().String(), timeoutHeight, 0)

	_, err = suite.chainC.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fungibleTokenPacket = types.NewFungibleTokenPacketData(fullDenomPathB2C, transferAmountDec.BigInt().String(), suite.chainC.SenderAccount().GetAddress().String(), suite.chainB.SenderAccount().GetAddress().String())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID, timeoutHeight, 0)
	err = pathBtoC.RelayPacket(packet, ack.Acknowledgement())
	suite.Require().NoError(err) // relay committed

	balanceB = suite.chainB.GetSimApp().BankKeeper.GetCoins(suite.chainB.GetContext(), suite.chainB.SenderAccount().GetAddress())
	// check that the balance on chainA returned back to the original state
	suite.Require().Equal(coinSentFromAToB, balanceB[0])

	// check that module account escrow address is empty
	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	balanceEC := suite.chainB.GetSimApp().BankKeeper.GetCoins(suite.chainB.GetContext(), escrowAddress)
	suite.Require().Nil(getBalance(balanceEC, denom))

	// check that balance on chain B is empty
	balanceC = suite.chainC.GetSimApp().BankKeeper.GetCoins(suite.chainC.GetContext(), suite.chainC.SenderAccount().GetAddress())
	suite.Require().Nil(getBalance(balanceC, denom))
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
