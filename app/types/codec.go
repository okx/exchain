package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/go-amino"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"math/big"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "okexchain/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)

	cdc.RegisterConcreteUnmarshaller(EthAccountName, func(data []byte) (n int, v interface{}, err error) {
		v, n, err = unmarshalEthAccountFromAmino(data)
		return
	})
}

func parsePosAndType(data byte) (pos int, aminoType amino.Typ3) {
	aminoType = amino.Typ3(data & 0x07)
	pos = int(data) >> 3
	return
}

func unmarshalEthAccountFromAmino(data []byte) (*EthAccount, int, error) {
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
			baseAccount, err := unmarshalBaseAccountFromAmino(subData)
			if err != nil {
				return nil, n, err
			}
			account.BaseAccount = baseAccount
		case 2:
			account.CodeHash = subData
		default:
			return nil, read, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return account, read, nil
}

func unmarshalBaseAccountFromAmino(data []byte) (*auth.BaseAccount, error) {
	var dataLen uint64 = 0
	var subData []byte
	account := &auth.BaseAccount{}

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, aminoType := parsePosAndType(data[0])
		data = data[1:]

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			account.Address = make([]byte, len(subData), len(subData))
			copy(account.Address, subData)
			// account.Address = subData
		case 2:
			coin, err := unmarshalCoinFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.Coins = append(account.Coins, coin)
		case 3:
			pubkey, err := unmarshalPubKeyFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.PubKey = pubkey
		case 4:
			uvarint, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			account.AccountNumber = uvarint
			dataLen = uint64(n)
		case 5:
			uvarint, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			account.Sequence = uvarint
			dataLen = uint64(n)
		}
	}
	return account, nil
}

func unmarshalCoinFromAmino(data []byte) (coin sdk.DecCoin, err error) {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, aminoType := parsePosAndType(data[0])
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
			amt := big.NewInt(0)
			err = amt.UnmarshalText(subData)
			if err != nil {
				return
			}
			coin.Amount = sdk.Dec{
				amt,
			}
		}
	}
	return
}

var typePubKeySecp256k1Prefix = []byte{0xeb, 0x5a, 0xe9, 0x87}

func unmarshalPubKeyFromAmino(data []byte) (tmcrypto.PubKey, error) {
	if 0 == bytes.Compare(typePubKeySecp256k1Prefix, data[0:4]) {
		if data[4] != 33 {
			return nil, errors.New("pubkey secp256k1 size error")
		}
		data = data[5:]
		pubKey := secp256k1.PubKeySecp256k1{}
		copy(pubKey[:], data)
		return pubKey, nil
	}
	panic("not implement")
}
