package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	cryptoAmino "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/codec"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil/testdata"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/legacy/legacytx"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
