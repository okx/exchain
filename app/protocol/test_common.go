package protocol

import (
	"encoding/hex"
	"os"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"testing"
)

func testPrepare(t *testing.T, db db.DB) (sdk.Context, storetypes.CommitMultiStore, proto.ProtocolKeeper) {
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(GetMainStoreKey(), sdk.StoreTypeIAVL, nil)
	require.Nil(t, ms.LoadLatestVersion())
	return sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout)), ms, GetEngine().GetProtocolKeeper()
}

func newPubKey(pubKey string) (res crypto.PubKey) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	var pubKeyEd25519 ed25519.PubKeyEd25519
	copy(pubKeyEd25519[:], pubKeyBytes[:])
	return pubKeyEd25519
}
