package ibc_tx

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/adapter"
	ibccodec "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/pb-codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx/internal/secp256k1"
)

var (
	PubKeyRegisterInterfaces = ibccodec.RegisterInterfaces
	LagacyKey2PbKey          = adapter.LagacyPubkey2ProtoBuffPubkey
	GenPrivKey               = secp256k1.GenPrivKey
)
