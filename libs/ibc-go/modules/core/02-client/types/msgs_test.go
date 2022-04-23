package types_test

import (
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	commitmenttypes "github.com/okex/exchain/libs/ibc-go/modules/core/23-commitment/types"
	solomachinetypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/06-solomachine/types"
	ibctmtypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

type TypesTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI
}

func (suite *TypesTestSuite) SetupTest() {
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

// tests that different clients within MsgCreateClient can be marshaled
// and unmarshaled.
func (suite *TypesTestSuite) TestMarshalMsgCreateClient() {
	var (
		msg *types.MsgCreateClient
		err error
	)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"solo machine client", func() {
			soloMachine := ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec(), "solomachine", "", 1)
			msg, err = types.NewMsgCreateClient(soloMachine.ClientState(), soloMachine.ConsensusState(), suite.chainA.SenderAccount().GetAddress())
			suite.Require().NoError(err)
		},
		},
		{
			"tendermint client", func() {
			tendermintClient := ibctmtypes.NewClientState(suite.chainA.ChainID(), ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
			msg, err = types.NewMsgCreateClient(tendermintClient, suite.chainA.CurrentTMClientHeader().ConsensusState(), suite.chainA.SenderAccount().GetAddress())
			suite.Require().NoError(err)
		},
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			tc.malleate()

			cdc := suite.chainA.App().AppCodec()

			// marshal message
			bz, err := cdc.GetProtocMarshal().MarshalJSON(msg)
			suite.Require().NoError(err)

			// unmarshal message
			newMsg := &types.MsgCreateClient{}
			err = cdc.GetProtocMarshal().UnmarshalJSON(bz, newMsg)
			suite.Require().NoError(err)

			suite.Require().True(proto.Equal(msg, newMsg))
		})
	}
}

// tests that different header within MsgUpdateClient can be marshaled
// and unmarshaled.

func (suite *TypesTestSuite) TestMarshalMsgUpgradeClient() {
	var (
		msg *types.MsgUpgradeClient
		err error
	)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"client upgrades to new tendermint client",
			func() {
				tendermintClient := ibctmtypes.NewClientState(suite.chainA.ChainID(), ibctesting.DefaultTrustLevel, ibctesting.TrustingPeriod, ibctesting.UnbondingPeriod, ibctesting.MaxClockDrift, clientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
				tendermintConsState := &ibctmtypes.ConsensusState{NextValidatorsHash: []byte("nextValsHash")}
				msg, err = types.NewMsgUpgradeClient("clientid", tendermintClient, tendermintConsState, []byte("proofUpgradeClient"), []byte("proofUpgradeConsState"), suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
		},
		{
			"client upgrades to new solomachine client",
			func() {
				soloMachine := ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec(), "solomachine", "", 1)
				msg, err = types.NewMsgUpgradeClient("clientid", soloMachine.ClientState(), soloMachine.ConsensusState(), []byte("proofUpgradeClient"), []byte("proofUpgradeConsState"), suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			tc.malleate()

			cdc := suite.chainA.App().AppCodec()

			// marshal message
			bz, err := cdc.GetProtocMarshal().MarshalJSON(msg)
			suite.Require().NoError(err)

			// unmarshal message
			newMsg := &types.MsgUpgradeClient{}
			err = cdc.GetProtocMarshal().UnmarshalJSON(bz, newMsg)
			suite.Require().NoError(err)
		})
	}
}

// tests that different misbehaviours within MsgSubmitMisbehaviour can be marshaled
// and unmarshaled.
func (suite *TypesTestSuite) TestMarshalMsgSubmitMisbehaviour() {
	var (
		msg *types.MsgSubmitMisbehaviour
		err error
	)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"solo machine client", func() {
			soloMachine := ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec(), "solomachine", "", 1)
			msg, err = types.NewMsgSubmitMisbehaviour(soloMachine.ClientID, soloMachine.CreateMisbehaviour(), suite.chainA.SenderAccount().GetAddress())
			suite.Require().NoError(err)
		},
		},
		{
			"tendermint client", func() {
			height := types.NewHeight(0, uint64(suite.chainA.CurrentHeader().Height))
			heightMinus1 := types.NewHeight(0, uint64(suite.chainA.CurrentHeader().Height)-1)
			header1 := suite.chainA.CreateTMClientHeader(suite.chainA.ChainID(), int64(height.RevisionHeight), heightMinus1, suite.chainA.CurrentHeader().Time, suite.chainA.Vals(), suite.chainA.Vals(), suite.chainA.Signers())
			header2 := suite.chainA.CreateTMClientHeader(suite.chainA.ChainID(), int64(height.RevisionHeight), heightMinus1, suite.chainA.CurrentHeader().Time.Add(time.Minute), suite.chainA.Vals(), suite.chainA.Vals(), suite.chainA.Signers())

			misbehaviour := ibctmtypes.NewMisbehaviour("tendermint", header1, header2)
			msg, err = types.NewMsgSubmitMisbehaviour("tendermint", misbehaviour, suite.chainA.SenderAccount().GetAddress())
			suite.Require().NoError(err)

		},
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			tc.malleate()

			cdc := suite.chainA.App().AppCodec()

			// marshal message
			bz, err := cdc.GetProtocMarshal().MarshalBinaryBare(msg)
			suite.Require().NoError(err)

			// unmarshal message
			newMsg := &types.MsgSubmitMisbehaviour{}
			err = cdc.GetProtocMarshal().UnmarshalBinaryBare(bz, newMsg)
			suite.Require().NoError(err)

			suite.Require().True(proto.Equal(msg, newMsg))
		})
	}
}

func (suite *TypesTestSuite) TestMsgSubmitMisbehaviour_ValidateBasic() {
	var (
		msg = &types.MsgSubmitMisbehaviour{}
		err error
	)

	cases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"invalid client-id",
			func() {
				msg.ClientId = ""
			},
			false,
		},
		{
			"valid - tendermint misbehaviour",
			func() {
				height := types.NewHeight(0, uint64(suite.chainA.CurrentHeader().Height))
				heightMinus1 := types.NewHeight(0, uint64(suite.chainA.CurrentHeader().Height)-1)
				header1 := suite.chainA.CreateTMClientHeader(suite.chainA.ChainID(), int64(height.RevisionHeight), heightMinus1, suite.chainA.CurrentHeader().Time, suite.chainA.Vals(), suite.chainA.Vals(), suite.chainA.Signers())
				header2 := suite.chainA.CreateTMClientHeader(suite.chainA.ChainID(), int64(height.RevisionHeight), heightMinus1, suite.chainA.CurrentHeader().Time.Add(time.Minute), suite.chainA.Vals(), suite.chainA.Vals(), suite.chainA.Signers())

				misbehaviour := ibctmtypes.NewMisbehaviour("tendermint", header1, header2)
				msg, err = types.NewMsgSubmitMisbehaviour("tendermint", misbehaviour, suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"invalid tendermint misbehaviour",
			func() {
				msg, err = types.NewMsgSubmitMisbehaviour("tendermint", &ibctmtypes.Misbehaviour{}, suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"failed to unpack misbehaviourt",
			func() {
				msg.Misbehaviour = nil
			},
			false,
		},
		{
			"invalid signer",
			func() {
				msg.Signer = ""
			},
			false,
		},
		{
			"valid - solomachine misbehaviour",
			func() {
				soloMachine := ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec(), "solomachine", "", 1)
				msg, err = types.NewMsgSubmitMisbehaviour(soloMachine.ClientID, soloMachine.CreateMisbehaviour(), suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"invalid solomachine misbehaviour",
			func() {
				msg, err = types.NewMsgSubmitMisbehaviour("solomachine", &solomachinetypes.Misbehaviour{}, suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"client-id mismatch",
			func() {
				soloMachineMisbehaviour := ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec(), "solomachine", "", 1).CreateMisbehaviour()
				msg, err = types.NewMsgSubmitMisbehaviour("external", soloMachineMisbehaviour, suite.chainA.SenderAccount().GetAddress())
				suite.Require().NoError(err)
			},
			false,
		},
	}

	for _, tc := range cases {
		tc.malleate()
		err = msg.ValidateBasic()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
