package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	auth "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/tendermint/go-amino"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "okexchain/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)

	cdc.RegisterConcreteUnmarshaller(EthAccountName, func(cdc *amino.Codec, data []byte) (v interface{}, n int, err error) {
		v, n, err = UnmarshalEthAccountFromAmino(cdc, data)
		return
	})
}

func UnmarshalEthAccountFromAmino(_ *amino.Codec, data []byte) (*EthAccount, int, error) {
	var dataLen uint64 = 0
	var read int
	var err error
	account := &EthAccount{}

	for {
		data = data[dataLen:]
		read += int(dataLen)

		if len(data) <= 0 {
			break
		}

		pos, _ := amino.ParseProtoPosAndTypeMustOneByte(data[0])
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
			baseAccount, err := auth.UnmarshalBaseAccountFromAmino(subData)
			if err != nil {
				return nil, n, err
			}
			account.BaseAccount = baseAccount
		case 2:
			account.CodeHash = make([]byte, len(subData))
			copy(account.CodeHash, subData)
		default:
			return nil, read, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return account, read, nil
}
