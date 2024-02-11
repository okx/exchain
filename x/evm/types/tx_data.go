package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/tendermint/go-amino"

	"github.com/okex/exchain/app/utils"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

// TxData implements the Ethereum transaction data structure. It is used
// solely as intended in Ethereum abiding by the protocol.
type TxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *ethcmn.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`

	// hash is only used when marshaling to JSON
	Hash *ethcmn.Hash `json:"hash" rlp:"-"`
}

// encodableTxData implements the Ethereum transaction data structure. It is used
// solely as intended in Ethereum abiding by the protocol.
type encodableTxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        string          `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *ethcmn.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       string          `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`

	// hash is only used when marshaling to JSON
	Hash *ethcmn.Hash `json:"hash" rlp:"-"`
}

func (tx *encodableTxData) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			data = data[n:]
			if len(data) < int(dataLen) {
				return fmt.Errorf("invalid tx data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			var n int
			tx.AccountNonce, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
		case 2:
			tx.Price = string(subData)
		case 3:
			var n int
			tx.GasLimit, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
		case 4:
			if dataLen != ethcmn.AddressLength {
				return errors.New("eth addr len error")
			}
			tx.Recipient = new(ethcmn.Address)
			copy(tx.Recipient[:], subData)
		case 5:
			tx.Amount = string(subData)
		case 6:
			tx.Payload = make([]byte, dataLen)
			copy(tx.Payload, subData)
		case 7:
			tx.V = string(subData)
		case 8:
			tx.R = string(subData)
		case 9:
			tx.S = string(subData)
		case 10:
			if dataLen != ethcmn.HashLength {
				return errors.New("hash len error")
			}
			tx.Hash = new(ethcmn.Hash)
			copy(tx.Hash[:], subData)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (td TxData) String() string {
	if td.Recipient != nil {
		return fmt.Sprintf("nonce=%d price=%s gasLimit=%d recipient=%s amount=%s data=0x%x v=%s r=%s s=%s",
			td.AccountNonce, td.Price, td.GasLimit, td.Recipient.Hex(), td.Amount, td.Payload, td.V, td.R, td.S)
	}

	return fmt.Sprintf("nonce=%d price=%s gasLimit=%d recipient=nil amount=%s data=0x%x v=%s r=%s s=%s",
		td.AccountNonce, td.Price, td.GasLimit, td.Amount, td.Payload, td.V, td.R, td.S)
}

// MarshalAmino defines custom encoding scheme for TxData
func (td TxData) MarshalAmino() ([]byte, error) {
	gasPrice, err := utils.MarshalBigInt(td.Price)
	if err != nil {
		return nil, err
	}

	amount, err := utils.MarshalBigInt(td.Amount)
	if err != nil {
		return nil, err
	}

	v, err := utils.MarshalBigInt(td.V)
	if err != nil {
		return nil, err
	}

	r, err := utils.MarshalBigInt(td.R)
	if err != nil {
		return nil, err
	}

	s, err := utils.MarshalBigInt(td.S)
	if err != nil {
		return nil, err
	}

	e := encodableTxData{
		AccountNonce: td.AccountNonce,
		Price:        gasPrice,
		GasLimit:     td.GasLimit,
		Recipient:    td.Recipient,
		Amount:       amount,
		Payload:      td.Payload,
		V:            v,
		R:            r,
		S:            s,
		Hash:         td.Hash,
	}

	return ModuleCdc.MarshalBinaryBare(e)
}

// UnmarshalAmino defines custom decoding scheme for TxData
func (td *TxData) UnmarshalAmino(data []byte) error {
	var e encodableTxData
	err := ModuleCdc.UnmarshalBinaryBare(data, &e)
	if err != nil {
		return err
	}

	td.AccountNonce = e.AccountNonce
	td.GasLimit = e.GasLimit
	td.Recipient = e.Recipient
	td.Payload = e.Payload
	td.Hash = e.Hash

	price, err := utils.UnmarshalBigInt(e.Price)
	if err != nil {
		return err
	}

	if td.Price != nil {
		td.Price.Set(price)
	} else {
		td.Price = price
	}

	amt, err := utils.UnmarshalBigInt(e.Amount)
	if err != nil {
		return err
	}

	if td.Amount != nil {
		td.Amount.Set(amt)
	} else {
		td.Amount = amt
	}

	v, err := utils.UnmarshalBigInt(e.V)
	if err != nil {
		return err
	}

	if td.V != nil {
		td.V.Set(v)
	} else {
		td.V = v
	}

	r, err := utils.UnmarshalBigInt(e.R)
	if err != nil {
		return err
	}

	if td.R != nil {
		td.R.Set(r)
	} else {
		td.R = r
	}

	s, err := utils.UnmarshalBigInt(e.S)
	if err != nil {
		return err
	}

	if td.S != nil {
		td.S.Set(s)
	} else {
		td.S = s
	}

	return nil
}

func (td *TxData) unmarshalFromAmino(cdc *amino.Codec, data []byte) error {
	var e encodableTxData
	err := e.UnmarshalFromAmino(cdc, data)
	if err != nil {
		return err
	}
	td.AccountNonce = e.AccountNonce
	td.GasLimit = e.GasLimit
	td.Recipient = e.Recipient
	td.Payload = e.Payload
	td.Hash = e.Hash

	price, err := utils.UnmarshalBigInt(e.Price)
	if err != nil {
		return err
	}

	if td.Price != nil {
		td.Price.Set(price)
	} else {
		td.Price = price
	}

	amt, err := utils.UnmarshalBigInt(e.Amount)
	if err != nil {
		return err
	}

	if td.Amount != nil {
		td.Amount.Set(amt)
	} else {
		td.Amount = amt
	}

	v, err := utils.UnmarshalBigInt(e.V)
	if err != nil {
		return err
	}

	if td.V != nil {
		td.V.Set(v)
	} else {
		td.V = v
	}

	r, err := utils.UnmarshalBigInt(e.R)
	if err != nil {
		return err
	}

	if td.R != nil {
		td.R.Set(r)
	} else {
		td.R = r
	}

	s, err := utils.UnmarshalBigInt(e.S)
	if err != nil {
		return err
	}

	if td.S != nil {
		td.S.Set(s)
	} else {
		td.S = s
	}

	return nil
}

func (td *TxData) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
	err := td.unmarshalFromAmino(cdc, data)
	if err != nil {
		u64, n, err := amino.DecodeUvarint(data)
		if err == nil && int(u64) == (len(data)-n) {
			return td.unmarshalFromAmino(cdc, data[n:])
		} else {
			return err
		}
	}
	return nil
}

func IsInscription(data []byte) bool {
	inscriptionStr := string(data)
	if strings.HasPrefix(inscriptionStr, "data:") {
		return true
	}
	first := strings.Index(inscriptionStr, "{")
	end := strings.Index(inscriptionStr, "}")

	if first < 0 || end < 0 || first >= end {
		return false
	}

	var obj interface{}
	err := json.Unmarshal([]byte(inscriptionStr[first:end+1]), &obj)
	if err == nil {
		return true
	}
	return false
}

// TODO: Implement JSON marshaling/ unmarshaling for this type

// TODO: Implement YAML marshaling/ unmarshaling for this type
