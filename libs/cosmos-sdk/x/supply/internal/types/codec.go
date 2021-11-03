package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/go-amino"
)

const (
	// MudulleAccountName is the amino encoding name for ModuleAccount
	MudulleAccountName = "cosmos-sdk/ModuleAccount"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, MudulleAccountName, nil)
	cdc.RegisterConcreteUnmarshaller(MudulleAccountName, func(cdc *amino.Codec, data []byte) (v interface{}, n int, err error) {
		v, n, err = UnmarshalMouduleAccountFromAmino(cdc, data)
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

func UnmarshalMouduleAccountFromAmino(_ *amino.Codec, data []byte) (*ModuleAccount, int, error) {
	var dataLen uint64 = 0
	var read int

	account := &ModuleAccount{}

	for {
		data = data[dataLen:]
		read += int(dataLen)

		if len(data) <= 0 {
			break
		}

		pos, _, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return nil, 0, err
		}
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
