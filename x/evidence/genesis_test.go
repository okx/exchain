package evidence_test

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/exchain/app"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/evidence"
	"github.com/okex/exchain/x/evidence/exported"
	"github.com/okex/exchain/x/evidence/internal/types"

	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper evidence.Keeper
}

func MakeOKEXApp() *app.ExChainApp {
	genesisState := app.NewDefaultGenesisState()
	db := dbm.NewMemDB()
	okexapp := app.NewExChainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)

	stateBytes, err := codec.MarshalJSONIndent(okexapp.Codec(), genesisState)
	if err != nil {
		panic(err)
	}
	okexapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	return okexapp
}

func (suite *GenesisTestSuite) SetupTest() {
	checkTx := false

	app := MakeOKEXApp()
	// get the app's codec and register custom testing types
	cdc := app.Codec()
	cdc.RegisterConcrete(types.TestEquivocationEvidence{}, "test/TestEquivocationEvidence", nil)

	// recreate keeper in order to use custom testing types
	evidenceKeeper := evidence.NewKeeper(
		cdc, app.GetKey(evidence.StoreKey), app.GetSubspace(evidence.ModuleName), app.StakingKeeper, app.SlashingKeeper,
	)
	router := evidence.NewRouter()
	router = router.AddRoute(types.TestEvidenceRouteEquivocation, types.TestEquivocationHandler(*evidenceKeeper))
	evidenceKeeper.SetRouter(router)

	suite.ctx = app.BaseApp.NewContext(checkTx, abci.Header{Height: 1})
	suite.keeper = *evidenceKeeper
}

func (suite *GenesisTestSuite) TestInitGenesis_Valid() {
	pk := ed25519.GenPrivKey()

	testEvidence := make([]exported.Evidence, 100)
	for i := 0; i < 100; i++ {
		sv := types.TestVote{
			ValidatorAddress: pk.PubKey().Address(),
			Height:           int64(i),
			Round:            0,
		}
		sig, err := pk.Sign(sv.SignBytes("test-chain"))
		suite.NoError(err)
		sv.Signature = sig

		testEvidence[i] = types.TestEquivocationEvidence{
			Power:      100,
			TotalPower: 100000,
			PubKey:     pk.PubKey(),
			VoteA:      sv,
			VoteB:      sv,
		}
	}

	suite.NotPanics(func() {
		evidence.InitGenesis(suite.ctx, suite.keeper, evidence.NewGenesisState(types.DefaultParams(), testEvidence))
	})

	for _, e := range testEvidence {
		_, ok := suite.keeper.GetEvidence(suite.ctx, e.Hash())
		suite.True(ok)
	}
}

func (suite *GenesisTestSuite) TestInitGenesis_Invalid() {
	pk := ed25519.GenPrivKey()

	testEvidence := make([]exported.Evidence, 100)
	for i := 0; i < 100; i++ {
		sv := types.TestVote{
			ValidatorAddress: pk.PubKey().Address(),
			Height:           int64(i),
			Round:            0,
		}
		sig, err := pk.Sign(sv.SignBytes("test-chain"))
		suite.NoError(err)
		sv.Signature = sig

		testEvidence[i] = types.TestEquivocationEvidence{
			Power:      100,
			TotalPower: 100000,
			PubKey:     pk.PubKey(),
			VoteA:      sv,
			VoteB:      types.TestVote{Height: 10, Round: 1},
		}
	}

	suite.Panics(func() {
		evidence.InitGenesis(suite.ctx, suite.keeper, evidence.NewGenesisState(types.DefaultParams(), testEvidence))
	})

	suite.Empty(suite.keeper.GetAllEvidence(suite.ctx))
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}
