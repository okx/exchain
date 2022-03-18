package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/tendermint/go-amino"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type KV struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// MarshalToAmino encode KV data to amino bytes
func (k *KV) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(k.Key) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, k.Key)
			if err != nil {
				return nil, err
			}

		case 2:
			if len(k.Value) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, k.Value)
			if err != nil {
				return nil, err
			}

		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal amino bytes to this object
func (k *KV) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			if len(data) < int(dataLen) {
				return errors.New("not enough data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			k.Key = make([]byte, len(subData))
			copy(k.Key, subData)

		case 2:
			k.Value = make([]byte, len(subData))
			copy(k.Value, subData)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
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
			dataLen, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}

			data = data[n:]
			if len(data) < int(dataLen) {
				return nil, fmt.Errorf("invalid data length: %d", dataLen)
			}
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

func (rd *ResultData) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
		dataLen, n, err = amino.DecodeUvarint(data)
		if err != nil {
			return err
		}

		data = data[n:]
		if len(data) < int(dataLen) {
			return errors.New("invalid data len")
		}
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
				return err
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

func (rd ResultData) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
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

	bz, err := data.MarshalToAmino(ModuleCdc)
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
		bz, err := amino.GetBinaryBareFromBinaryLengthPrefixed(in)
		if err == nil {
			var data ResultData
			err = data.UnmarshalFromAmino(ModuleCdc, bz)
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

var ethAddrStringPool = &sync.Pool{
	New: func() interface{} {
		return &[32]byte{}
	},
}

type EthAddressStringer ethcmn.Address

func (address EthAddressStringer) String() string {
	p := &address
	return EthAddressToString((*ethcmn.Address)(p))
}

func EthAddressToString(address *ethcmn.Address) string {
	var buf [len(address)*2 + 2]byte
	buf[0] = '0'
	buf[1] = 'x'
	hex.Encode(buf[2:], address[:])

	// compute checksum
	sha := keccakStatePool.Get().(ethcrypto.KeccakState)
	sha.Reset()
	sha.Write(buf[2:])

	hash := ethAddrStringPool.Get().(*[32]byte)
	sha.Read(hash[:])

	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	ethAddrStringPool.Put(hash)
	keccakStatePool.Put(sha)
	return amino.BytesToStr(buf[:])
}

type EthHashStringer ethcmn.Hash

func (h EthHashStringer) String() string {
	var enc [len(h)*2 + 2]byte
	enc[0] = '0'
	enc[1] = 'x'
	hex.Encode(enc[2:], h[:])
	return amino.BytesToStr(enc[:])
}
