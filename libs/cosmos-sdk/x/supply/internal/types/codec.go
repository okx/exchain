package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/go-amino"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "cosmos-sdk/ModuleAccount", nil)
	cdc.RegisterConcreteUnmarshaller("cosmos-sdk/ModuleAccount", func(data []byte) (n int, v interface{}, err error) {
		v, n, err = unmarshalMouduleAccountFromAmino(data)
		return
	})

	cdc.RegisterConcrete(&Supply{}, "cosmos-sdk/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}

func parsePosAndType(data byte) (pos int, aminoType amino.Typ3) {
	aminoType = amino.Typ3(data & 0x07)
	pos = int(data) >> 3
	return
}

func unmarshalMouduleAccountFromAmino(data []byte) (*ModuleAccount, int, error) {
	var dataLen uint64 = 0
	var read int
	var err error
	account := &ModuleAccount{}

	for {
		data = data[dataLen:]
		read += int(dataLen)

		if len(data) <= 0 {
			break
		}

		pos, _ := parsePosAndType(data[0])
		data = data[1:]
		read += 1

		var n int
		dataLen, n, err = amino.DecodeUvarint(data)
		if err != nil {
			return nil, read, err
		}

		data = data[n:]
		read += n
		subData := data[:dataLen]

		switch pos {
		case 1:
			baseAccount, err := types.UnmarshalBaseAccountFromAmino(subData)
			if err != nil {
				return nil, n, err
			}
			account.BaseAccount = baseAccount
		case 2:
			account.Name = string(subData)
		case 3:
			account.Permissions = append(account.Permissions, string(subData))
		default:
			return nil, read, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return account, read, nil
}
