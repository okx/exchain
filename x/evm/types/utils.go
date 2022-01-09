package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"math/big"
	"strings"

	"github.com/tendermint/go-amino"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type KV struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// GenerateEthAddress generates an Ethereum address.
func GenerateEthAddress() ethcmn.Address {
	priv, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	return ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
}

// ValidateSigner attempts to validate a signer for a given slice of bytes over
// which a signature and signer is given. An error is returned if address
// derived from the signature and bytes signed does not match the given signer.
func ValidateSigner(signBytes, sig []byte, signer ethcmn.Address) error {
	pk, err := ethcrypto.SigToPub(signBytes, sig)

	if err != nil {
		return errors.Wrap(err, "failed to derive public key from signature")
	} else if ethcrypto.PubkeyToAddress(*pk) != signer {
		return fmt.Errorf("invalid signature for signer: %s", signer)
	}

	return nil
}

func rlpHash(x interface{}) (hash ethcmn.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	_ = rlp.Encode(hasher, x)
	_ = hasher.Sum(hash[:0])

	return hash
}

// ResultData represents the data returned in an sdk.Result
type ResultData struct {
	ContractAddress ethcmn.Address  `json:"contract_address"`
	Bloom           ethtypes.Bloom  `json:"bloom"`
	Logs            []*ethtypes.Log `json:"logs"`
	Ret             []byte          `json:"ret"`
	TxHash          ethcmn.Hash     `json:"tx_hash"`
}

func UnmarshalEthLogFromAmino(data []byte) (*ethtypes.Log, error) {
	var dataLen uint64 = 0
	var subData []byte
	log := &ethtypes.Log{}

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
			if int(dataLen) != ethcmn.AddressLength {
				return nil, fmt.Errorf("invalid address length: %d", dataLen)
			}
			copy(log.Address[:], subData)
		case 2:
			if int(dataLen) != ethcmn.HashLength {
				return nil, fmt.Errorf("invalid topic hash length: %d", dataLen)
			}
			var hash ethcmn.Hash
			copy(hash[:], subData)
			log.Topics = append(log.Topics, hash)
		case 3:
			log.Data = make([]byte, dataLen)
			copy(log.Data, subData)
		case 4:
			var n int
			log.BlockNumber, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			dataLen = uint64(n)
		case 5:
			if int(dataLen) != ethcmn.HashLength {
				return nil, fmt.Errorf("invalid topic hash length: %d", dataLen)
			}
			copy(log.TxHash[:], subData)
		case 6:
			var n int
			var uv uint64
			uv, n, err = amino.DecodeUvarint(data)
			log.TxIndex = uint(uv)
			if err != nil {
				return nil, err
			}
			dataLen = uint64(n)
		case 7:
			if int(dataLen) != ethcmn.HashLength {
				return nil, fmt.Errorf("invalid topic hash length: %d", dataLen)
			}
			copy(log.BlockHash[:], subData)
		case 8:
			var n int
			var uv uint64
			uv, n, err = amino.DecodeUvarint(data)
			log.Index = uint(uv)
			if err != nil {
				return nil, err
			}
			dataLen = uint64(n)
		case 9:
			if data[0] == 0 {
				log.Removed = false
			} else if data[0] == 1 {
				log.Removed = true
			} else {
				return nil, fmt.Errorf("invalid removed flag: %d", data[0])
			}
			dataLen = 1
		}
	}
	return log, nil
}

var ethLogBufferPool = amino.NewBufferPool()

func MarshalEthLogToAmino(log *ethtypes.Log) ([]byte, error) {
	if log == nil {
		return nil, nil
	}
	var buf = ethLogBufferPool.Get()
	defer ethLogBufferPool.Put(buf)
	fieldKeysType := [9]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2, 4 << 3, 5<<3 | 2, 6 << 3, 7<<3 | 2, 8 << 3, 9 << 3}
	for pos := 1; pos < 10; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err := buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}

		switch pos {
		case 1:
			err := buf.WriteByte(byte(ethcmn.AddressLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(log.Address.Bytes())
			if err != nil {
				return nil, err
			}
		case 2:
			topicsLen := len(log.Topics)
			if topicsLen == 0 {
				noWrite = true
				break
			}
			err = buf.WriteByte(byte(ethcmn.HashLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(log.Topics[0].Bytes())
			if err != nil {
				return nil, err
			}

			for i := 1; i < topicsLen; i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}

				err = buf.WriteByte(byte(ethcmn.HashLength))
				if err != nil {
					return nil, err
				}
				_, err = buf.Write(log.Topics[i].Bytes())
				if err != nil {
					return nil, err
				}
			}
		case 3:
			dataLen := len(log.Data)
			if dataLen == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(buf, uint64(dataLen))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(log.Data)
			if err != nil {
				return nil, err
			}
		case 4:
			if log.BlockNumber == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(buf, log.BlockNumber)
			if err != nil {
				return nil, err
			}
		case 5:
			err := buf.WriteByte(byte(ethcmn.HashLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(log.TxHash.Bytes())
			if err != nil {
				return nil, err
			}
		case 6:
			if log.TxIndex == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarintToBuffer(buf, uint64(log.TxIndex))
			if err != nil {
				return nil, err
			}
		case 7:
			err := buf.WriteByte(byte(ethcmn.HashLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(log.BlockHash.Bytes())
			if err != nil {
				return nil, err
			}
		case 8:
			if log.Index == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarintToBuffer(buf, uint64(log.Index))
			if err != nil {
				return nil, err
			}
		case 9:
			if log.Removed {
				err = buf.WriteByte(byte(1))
				if err != nil {
					return nil, err
				}
			} else {
				noWrite = true
				break
			}
		default:
			panic("unreachable")
		}

		if noWrite {
			buf.Truncate(lBeforeKey)
		}
	}
	return amino.GetBytesBufferCopy(buf), nil
}

func (rd *ResultData) UnmarshalFromAmino(data []byte) error {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		if aminoType != amino.Typ3_ByteLength {
			return fmt.Errorf("unexpect proto type %d", aminoType)
		}
		data = data[1:]

		var n int
		dataLen, n, _ = amino.DecodeUvarint(data)

		data = data[n:]
		subData = data[:dataLen]

		switch pos {
		case 1:
			if int(dataLen) != ethcmn.AddressLength {
				return fmt.Errorf("invalid contract address length: %d", dataLen)
			}
			copy(rd.ContractAddress[:], subData)
		case 2:
			if int(dataLen) != ethtypes.BloomByteLength {
				return fmt.Errorf("invalid bloom length: %d", dataLen)
			}
			copy(rd.Bloom[:], subData)
		case 3:
			var log *ethtypes.Log
			if dataLen == 0 {
				log, err = nil, nil
			} else {
				log, err = UnmarshalEthLogFromAmino(subData)
			}
			if err != nil {
				return nil
			}
			rd.Logs = append(rd.Logs, log)
		case 4:
			rd.Ret = make([]byte, dataLen)
			copy(rd.Ret, subData)
		case 5:
			if dataLen != ethcmn.HashLength {
				return fmt.Errorf("invalid tx hash length %d", dataLen)
			}
			copy(rd.TxHash[:], subData)
		}
	}
	return nil
}

var resultDataBufferPool = amino.NewBufferPool()

func (rd ResultData) MarshalToAmino() ([]byte, error) {
	var buf = resultDataBufferPool.Get()
	defer resultDataBufferPool.Put(buf)
	fieldKeysType := [5]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2, 4<<3 | 2, 5<<3 | 2}
	for pos := 1; pos < 6; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err := buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}

		switch pos {
		case 1:
			err := buf.WriteByte(byte(ethcmn.AddressLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(rd.ContractAddress.Bytes())
			if err != nil {
				return nil, err
			}
		case 2:
			_, err := buf.Write([]byte{0b10000000, 0b00000010}) // bloom length 256
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(rd.Bloom.Bytes())
			if err != nil {
				return nil, err
			}
		case 3:
			logsLen := len(rd.Logs)
			if logsLen == 0 {
				noWrite = true
				break
			}
			data, err := MarshalEthLogToAmino(rd.Logs[0])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(buf, uint64(len(data)))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(data)
			if err != nil {
				return nil, err
			}
			for i := 1; i < logsLen; i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err = MarshalEthLogToAmino(rd.Logs[i])
				if err != nil {
					return nil, err
				}
				err = amino.EncodeUvarintToBuffer(buf, uint64(len(data)))
				if err != nil {
					return nil, err
				}
				_, err = buf.Write(data)
				if err != nil {
					return nil, err
				}
			}
		case 4:
			retLen := len(rd.Ret)
			if retLen == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarintToBuffer(buf, uint64(retLen))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(rd.Ret)
			if err != nil {
				return nil, err
			}
		case 5:
			err := buf.WriteByte(byte(ethcmn.HashLength))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(rd.TxHash.Bytes())
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}

		if noWrite {
			buf.Truncate(lBeforeKey)
		}
	}
	return amino.GetBytesBufferCopy(buf), nil
}

// String implements fmt.Stringer interface.
func (rd ResultData) String() string {
	var logsStr string
	logsLen := len(rd.Logs)
	for i := 0; i < logsLen; i++ {
		logsStr = fmt.Sprintf("%s\t\t%v\n ", logsStr, *rd.Logs[i])
	}

	return strings.TrimSpace(fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %s
	Ret: %v
	TxHash: %s	
	Logs: 
%s`, rd.ContractAddress.String(), rd.Bloom.Big().String(), rd.Ret, rd.TxHash.String(), logsStr))
}

// EncodeResultData takes all of the necessary data from the EVM execution
// and returns the data as a byte slice encoded with amino
func EncodeResultData(data ResultData) ([]byte, error) {
	var buf = new(bytes.Buffer)

	bz, err := data.MarshalToAmino()
	if err != nil {
		bz, err = ModuleCdc.MarshalBinaryBare(data)
		if err != nil {
			return nil, err
		}
	}

	// Write uvarint(len(bz)).
	err = amino.EncodeUvarintToBuffer(buf, uint64(len(bz)))
	if err != nil {
		return nil, err
	}

	// Write bz.
	_, err = buf.Write(bz)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeResultData decodes an amino-encoded byte slice into ResultData
func DecodeResultData(in []byte) (ResultData, error) {
	if len(in) > 0 {
		u64, n := binary.Uvarint(in)
		if u64 == uint64(len(in)-n) {
			bz := in[n:]
			var data ResultData
			err := data.UnmarshalFromAmino(bz)
			if err == nil {
				return data, nil
			}
		}
	}
	var data ResultData
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(in, &data)
	if err != nil {
		return ResultData{}, err
	}
	return data, nil
}

// ----------------------------------------------------------------------------
// Auxiliary

type decoderConfig struct {
	enableRLP bool // disable RLP by default
}

// Option is function to choose decoder.
type Option func(config *decoderConfig)

// EnableRLP enables RLP decoder.
func EnableRLP() Option {
	return func(config *decoderConfig) {
		config.enableRLP = true
	}
}

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec, options ...Option) sdk.TxDecoder {
	var cfg decoderConfig
	for _, op := range options {
		op(&cfg)
	}
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx sdk.Tx
		var err error
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		if cfg.enableRLP {
			// Try to decode as MsgEthereumTx through RLP
			var ethTx MsgEthereumTx
			if err = authtypes.EthereumTxDecode(txBytes, &ethTx); err == nil {
				return ethTx, nil
			}
		}

		// sdk.Tx is an interface. The concrete message types
		// are registered by MakeTxCodec
		// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
		if v, err := cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err == nil {
			return v.(sdk.Tx), nil
		}
		if err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err == nil {
			return tx, nil
		}
		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// recoverEthSig recovers a signature according to the Ethereum specification and
// returns the sender or an error.
//
// Ref: Ethereum Yellow Paper (BYZANTIUM VERSION 69351d5) Appendix F
// nolint: gocritic
func recoverEthSig(R, S, Vb *big.Int, sigHash ethcmn.Hash) (ethcmn.Address, error) {
	if Vb.BitLen() > 8 {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	V := byte(Vb.Uint64() - 27)
	if !ethcrypto.ValidateSignatureValues(V, R, S, true) {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)

	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	pub, err := ethcrypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}

	var addr ethcmn.Address
	copy(addr[:], ethcrypto.Keccak256(pub[1:])[12:])

	return addr, nil
}
