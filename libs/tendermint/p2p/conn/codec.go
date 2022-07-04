package conn

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
)

var cdc *amino.Codec = amino.NewCodec()

var packetMsgAminoTypePrefix []byte

func init() {
	cryptoamino.RegisterAmino(cdc)
	RegisterPacket(cdc)

	packetMsgAminoTypePrefix = initPacketMsgAminoTypePrefix(cdc)
}

func initPacketMsgAminoTypePrefix(cdc *amino.Codec) []byte {
	packetMsgAminoTypePrefix := make([]byte, 8)
	tpl, err := cdc.GetTypePrefix(PacketMsg{}, packetMsgAminoTypePrefix)
	if err != nil {
		panic(err)
	}
	packetMsgAminoTypePrefix = packetMsgAminoTypePrefix[:tpl]
	return packetMsgAminoTypePrefix
}

func getPacketMsgAminoTypePrefix() []byte {
	return packetMsgAminoTypePrefix
}
