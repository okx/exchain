package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
)

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}

func UnmarshalCoinFromAmino(data []byte) (coin DecCoin, err error) {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return coin, err
		}
		data = data[1:]

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			coin.Denom = string(subData)
		case 2:
			err = coin.Amount.UnmarshalFromAmino(subData)
			if err != nil {
				return coin, err
			}
		}
	}
	return coin, err
}
