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
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/sr25519"
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
var typePubKeyEd25519Prefix = []byte{0x16, 0x24, 0xde, 0x64}
var typePubKeySr25519Prefix = []byte{0x0d, 0xfb, 0x10, 0x05}

func unmarshalPubKeyFromAmino(data []byte) (tmcrypto.PubKey, error) {
	if data[0] == 0x00 {
		return nil, errors.New("unmarshal pubkey with disamb do not implement")
	}
	prefix := data[0:4]
	size := data[4]
	data = data[5:]
	if len(data) < int(size) {
		return nil, errors.New("raw data size error")
	}
	if 0 == bytes.Compare(typePubKeySecp256k1Prefix, prefix) {
		if size != secp256k1.PubKeySecp256k1Size {
			return nil, errors.New("pubkey secp256k1 size error")
		}
		pubKey := secp256k1.PubKeySecp256k1{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else if 0 == bytes.Compare(typePubKeyEd25519Prefix, prefix) {
		if size != ed25519.PubKeyEd25519Size {
			return nil, errors.New("pubkey ed25519 size error")
		}
		pubKey := ed25519.PubKeyEd25519{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else if 0 == bytes.Compare(typePubKeySr25519Prefix, prefix) {
		if size != sr25519.PubKeySr25519Size {
			return nil, errors.New("pubkey sr25519 size error")
		}
		pubKey := sr25519.PubKeySr25519{}
		copy(pubKey[:], data)
		return pubKey, nil
	} else {
		return nil, errors.New("unmarshal pubkey do not implement")
	}
}
