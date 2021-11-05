package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	cryptoamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/tendermint/go-amino"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil)

	cdc.RegisterConcreteUnmarshaller("cosmos-sdk/StdTx", func(c *amino.Codec, bytes []byte) (interface{}, int, error) {
		var tx StdTx
		err := tx.UnmarshalFromAmino(c, bytes)
		if err != nil {
			return nil, 0, err
		}
		return tx, len(bytes), nil
	})
}

// RegisterAccountTypeCodec registers an external account type defined in
// another module for the internal ModuleCdc.
func RegisterAccountTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}

func UnmarshalBaseAccountFromAmino(data []byte) (*BaseAccount, error) {
	var dataLen uint64 = 0
	var subData []byte
	account := &BaseAccount{}

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return nil, err
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
			account.Address = make([]byte, len(subData))
			copy(account.Address, subData)
			// account.Address = subData
		case 2:
			coin, err := sdk.UnmarshalCoinFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.Coins = append(account.Coins, coin)
		case 3:
			pubkey, err := cryptoamino.UnmarshalPubKeyFromAminoWithTypePrefix(subData)
			if err != nil {
				return nil, err
			}
			account.PubKey = pubkey
		case 4:
			var n int
			account.AccountNumber, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			dataLen = uint64(n)
		case 5:
			var n int
			account.Sequence, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			dataLen = uint64(n)
		}
	}
	return account, nil
}
