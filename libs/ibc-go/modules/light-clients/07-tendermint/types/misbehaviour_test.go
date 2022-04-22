package types_test

import (
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
	"github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
)

func (suite *TendermintTestSuite) TestMisbehaviour() {

	signers := []tmtypes.PrivValidator{suite.privVal}
	heightMinus1 := clienttypes.NewHeight(0, height.RevisionHeight-1)

	misbehaviour := &types.Misbehaviour{
		Header1:  suite.header,
		Header2:  suite.chainA.CreateTMClientHeader(chainID, int64(height.RevisionHeight), heightMinus1, suite.now, suite.valSet, suite.valSet, signers),
		ClientId: clientID,
	}

	suite.Require().Equal(exported.Tendermint, misbehaviour.ClientType())
	suite.Require().Equal(clientID, misbehaviour.GetClientID())
}
