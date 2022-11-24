package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

var (
	_ exported.SignerPacketI = SignerPacketWrapper{}
)

type SignerPacketWrapper struct {
	Packet
	Signers []sdk.AccAddress
	Gas     sdk.Gas
}

func NewSignerPacketWrapper(packet Packet, signers []sdk.AccAddress, gas sdk.Gas) *SignerPacketWrapper {
	return &SignerPacketWrapper{Packet: packet, Signers: signers, Gas: gas}
}

func (s SignerPacketWrapper) GetSigner() []sdk.AccAddress {
	return s.Signers
}

func (s SignerPacketWrapper) GetInternal() exported.PacketI {
	return s.Packet
}

func (s SignerPacketWrapper) GetGas() sdk.Gas {
	return s.Gas
}
