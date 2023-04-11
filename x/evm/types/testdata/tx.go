package testdata

import (
	"encoding/json"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// KeyTestPubAddr generates a new secp256k1 keypair.
//func KeyTestPubAddr() (cryptotypes.PrivKey, cryptotypes.PubKey, sdk.AccAddress) {
//	key := secp256k1.GenPrivKey()
//	pub := key.PubKey()
//	addr := sdk.AccAddress(pub.Address())
//	return key, pub, addr
//}

// KeyTestPubAddr generates a new secp256r1 keypair.
//func KeyTestPubAddrSecp256R1(require *require.Assertions) (cryptotypes.PrivKey, cryptotypes.PubKey, sdk.AccAddress) {
//	key, err := secp256r1.GenPrivKey()
//	require.NoError(err)
//	pub := key.PubKey()
//	addr := sdk.AccAddress(pub.Address())
//	return key, pub, addr
//}


var _ sdk.Msg = (*TestMsg)(nil)

func (msg *TestMsg) Route() string { return "TestMsg" }
func (msg *TestMsg) Type() string  { return "Test message" }
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.Signers)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}
func (msg *TestMsg) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Signers))
	for i, in := range msg.Signers {
		addr, err := sdk.AccAddressFromBech32(in)
		if err != nil {
			panic(err)
		}

		addrs[i] = addr
	}

	return addrs
}
func (msg *TestMsg) ValidateBasic() error { return nil }

var _ ibcmsg.Msg = &MsgCreateDog{}

func (msg *MsgCreateDog) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{} }
func (msg *MsgCreateDog) ValidateBasic() error         { return nil }
