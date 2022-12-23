package keeper_test

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/okex/exchain/libs/ibc-go/modules/core/23-commitment/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	ibctmtypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	types2 "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/x/erc20/keeper"
)

const CorrectIbcDenom2 = "ibc/3EF3B49764DB0E2284467F8BF7A08C18EACACB30E1AD7ABA8E892F1F679443F9"
const CorrectIbcDenom = "ibc/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

func (suite *KeeperTestSuite) TestSendToIbcHandler() {
	contract := common.BigToAddress(big.NewInt(1))
	sender := common.BigToAddress(big.NewInt(2))
	invalidDenom := "testdenom"
	validDenom := CorrectIbcDenom
	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"non associated coin denom, expect fail",
			func() {
				coin := sdk.NewCoin(invalidDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), invalidDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			errors.New("contract 0x0000000000000000000000000000000000000001 is not connected to native token"),
		},
		{
			"non IBC denom, expect fail",
			func() {
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, invalidDenom, contract)
				coin := sdk.NewCoin(invalidDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), invalidDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			errors.New("the native token associated with the contract 0x0000000000000000000000000000000000000001 is not an ibc voucher"),
		},
		{
			"success send to ibc",
			func() {
				amount := sdk.NewInt(100)
				suite.app.TransferKeeper.SetParams(suite.ctx, types2.Params{
					true, true,
				})
				channelA := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, types2.DenomTrace{
					BaseDenom: "ibc/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", Path: "",
				})
				suite.app.TransferKeeper.BindPort(suite.ctx, "transfer")
				cap, _ := suite.app.ScopedTransferKeeper.NewCapability(suite.ctx, host.ChannelCapabilityPath("transfer", channelA))
				suite.app.ScopedIBCKeeper.ClaimCapability(suite.ctx, cap, host.ChannelCapabilityPath("transfer", channelA))
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, CorrectIbcDenom2, contract)
				c := channeltypes.Channel{
					State:    channeltypes.OPEN,
					Ordering: channeltypes.UNORDERED,
					Counterparty: channeltypes.Counterparty{
						PortId:    "transfer",
						ChannelId: channelA,
					},
					ConnectionHops: []string{"one"},
					Version:        "version",
				}

				suite.app.IBCKeeper.V2Keeper.ChannelKeeper.SetNextSequenceSend(suite.ctx, "transfer", channelA, 1)
				suite.app.IBCKeeper.V2Keeper.ChannelKeeper.SetChannel(suite.ctx, "transfer", channelA, c)
				counterparty := connectiontypes.NewCounterparty("client-1", "one", commitmenttypes.NewMerklePrefix([]byte("ibc")))
				conn1 := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, "client-1", counterparty, connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions()), 0)
				period := time.Hour * 24 * 7 * 2
				clientState := ibctmtypes.NewClientState("testChainID", ibctmtypes.DefaultTrustLevel, period, period, period, types.NewHeight(0, 5), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
				suite.app.IBCKeeper.V2Keeper.ClientKeeper.SetClientState(suite.ctx, "client-1", clientState)
				consensusState := ibctmtypes.NewConsensusState(time.Now(), commitmenttypes.NewMerkleRoot([]byte("root")), []byte("nextValsHash"))
				suite.app.IBCKeeper.V2Keeper.ClientKeeper.SetClientConsensusState(suite.ctx, "client-1", types.NewHeight(0, 5), consensusState)
				suite.app.IBCKeeper.V2Keeper.ConnectionKeeper.SetConnection(suite.ctx, "one", conn1)
				coin := sdk.NewCoin(CorrectIbcDenom2, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), CorrectIbcDenom2)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			nil,
		},
		{
			"denomination trace not found",
			func() {
				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			errors.New("denomination trace not found: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			handler := keeper.NewSendToIbcEventHandler(suite.app.Erc20Keeper)
			tc.malleate()
			err := handler.Handle(suite.ctx, contract, data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendNative20ToIbcHandler() {
	contract := common.BigToAddress(big.NewInt(1))
	sender := common.BigToAddress(big.NewInt(2))
	validDenom := "testdenom"

	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"non associated coin denom, expect fail",
			func() {
				coin := sdk.NewCoin(validDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			errors.New("contract 0x0000000000000000000000000000000000000001 is not connected to native token"),
		},
		{
			"success send to ibc",
			func() {
				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					amount.BigInt(),
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			nil,
		},
		{
			"portid channel error",
			func() {

				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)
				suite.ctx.SetIsCheckTx(false)
				suite.ctx.SetIsTraceTx(false)
				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)
				suite.app.TransferKeeper.SetParams(suite.ctx, types2.Params{true, true})
				input, err := keeper.SendNative20ToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
					"transfer",
					"channel-0",
				)
				suite.Require().NoError(err)
				data = input
			},
			func() {},
			errors.New("channel not found: port ID (transfer) channel ID (channel-0)"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			handler := keeper.NewSendNative20ToIbcEventHandler(suite.app.Erc20Keeper)
			tc.malleate()
			err := handler.Handle(suite.ctx, contract, data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}
