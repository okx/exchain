package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/okex/exchain/libs/tendermint/libs/kv"

	"github.com/tendermint/go-amino"
)

func (pubkey PubKey) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if pubkey.Type == "" {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeStringToBuffer(&buf, pubkey.Type)
			if err != nil {
				return nil, err
			}
		case 2:
			if len(pubkey.Data) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, pubkey.Data)
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

func (valUpdate ValidatorUpdate) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [2]byte{1<<3 | 2, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool

		switch pos {
		case 1:
			var data []byte
			data, err = valUpdate.PubKey.MarshalToAmino()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				noWrite = true
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, data)
			if err != nil {
				return nil, err
			}
		case 2:
			if valUpdate.Power == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(valUpdate.Power))
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
	return buf.Bytes(), nil
}

func (params BlockParams) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1 << 3, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err := buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}
		switch pos {
		case 1:
			if params.MaxBytes == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(params.MaxBytes))
			if err != nil {
				return nil, err
			}
		case 2:
			if params.MaxGas == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(params.MaxGas))
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
	return buf.Bytes(), nil
}

func (params EvidenceParams) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1 << 3, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err := buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}
		switch pos {
		case 1:
			if params.MaxAgeNumBlocks == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(params.MaxAgeNumBlocks))
			if err != nil {
				return nil, err
			}
		case 2:
			if params.MaxAgeDuration == 0 {
				noWrite = true
				break
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(params.MaxAgeDuration))
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
	return buf.Bytes(), nil
}

func (params ValidatorParams) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [1]byte{1<<3 | 2}
	for pos := 1; pos <= 1; pos++ {
		switch pos {
		case 1:
			if len(params.PubKeyTypes) == 0 {
				break
			}
			for i := 0; i < len(params.PubKeyTypes); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				err = amino.EncodeStringToBuffer(&buf, params.PubKeyTypes[i])
				if err != nil {
					return nil, err
				}
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

func (event Event) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if event.Type == "" {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(len(event.Type)))
			if err != nil {
				return nil, err
			}
			_, err = buf.WriteString(event.Type)
			if err != nil {
				return nil, err
			}
		case 2:
			for i := 0; i < len(event.Attributes); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := kv.MarshalPairToAmino(event.Attributes[i])
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

func (event *Event) UnmarshalFromAmino(data []byte) error {
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
			return errors.New("invalid amino type")
		}
		data = data[1:]

		var n int
		dataLen, n, _ = amino.DecodeUvarint(data)

		data = data[n:]
		subData = data[:dataLen]

		switch pos {
		case 1:
			event.Type = string(subData)
		case 2:
			var kvpair kv.Pair
			err = kvpair.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			event.Attributes = append(event.Attributes, kvpair)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (tx *ResponseDeliverTx) MarshalToAmino() ([]byte, error) {
	if tx == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	var err error
	fieldKeysType := [8]byte{1 << 3, 2<<3 | 2, 3<<3 | 2, 4<<3 | 2, 5 << 3, 6 << 3, 7<<3 | 2, 8<<3 | 2}
	for pos := 1; pos <= 8; pos++ {
		switch pos {
		case 1:
			if tx.Code == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(tx.Code))
			if err != nil {
				return nil, err
			}
		case 2:
			if len(tx.Data) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, tx.Data)
			if err != nil {
				return nil, err
			}
		case 3:
			if tx.Log == "" {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeStringToBuffer(&buf, tx.Log)
			if err != nil {
				return nil, err
			}
		case 4:
			if tx.Info == "" {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeStringToBuffer(&buf, tx.Info)
			if err != nil {
				return nil, err
			}
		case 5:
			if tx.GasWanted == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(tx.GasWanted))
			if err != nil {
				return nil, err
			}
		case 6:
			if tx.GasUsed == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(tx.GasUsed))
			if err != nil {
				return nil, err
			}
		case 7:
			for i := 0; i < len(tx.Events); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := tx.Events[i].MarshalToAmino()
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		case 8:
			if tx.Codespace == "" {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeStringToBuffer(&buf, tx.Codespace)
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

func (tx *ResponseDeliverTx) UnmarshalFromAmino(data []byte) error {
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
		data = data[1:]

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			tx.Code = uint32(uvint)
			dataLen = uint64(n)
		case 2:
			tx.Data = make([]byte, dataLen)
			copy(tx.Data, subData)
		case 3:
			tx.Log = string(subData)
		case 4:
			tx.Info = string(subData)
		case 5:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			tx.GasWanted = int64(uvint)
			dataLen = uint64(n)
		case 6:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			tx.GasUsed = int64(uvint)
			dataLen = uint64(n)
		case 7:
			var event Event
			err = event.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			tx.Events = append(tx.Events, event)
		case 8:
			tx.Codespace = string(subData)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (beginBlock ResponseBeginBlock) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKey := byte(1<<3 | 2)
	for i := 0; i < len(beginBlock.Events); i++ {
		err := buf.WriteByte(fieldKey)
		if err != nil {
			return nil, err
		}
		data, err := beginBlock.Events[i].MarshalToAmino()
		if err != nil {
			return nil, err
		}
		err = amino.EncodeByteSliceToBuffer(&buf, data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func MarshalConsensusParamsToAmino(params ConsensusParams) (data []byte, err error) {
	var buf bytes.Buffer
	fieldKeysType := [3]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2}
	for pos := 1; pos <= 3; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err = buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}
		switch pos {
		case 1:
			if params.Block == nil {
				noWrite = true
				break
			}
			data, err = params.Block.MarshalToAmino()
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, data)
			if err != nil {
				return nil, err
			}
		case 2:
			if params.Evidence == nil {
				noWrite = true
				break
			}
			data, err = params.Evidence.MarshalToAmino()
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, data)
			if err != nil {
				return nil, err
			}
		case 3:
			if params.Validator == nil {
				noWrite = true
				break
			}
			data, err = params.Validator.MarshalToAmino()
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, data)
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
	return buf.Bytes(), nil
}

func MarshalResponseEndBlockToAmino(endBlock *ResponseEndBlock) ([]byte, error) {
	if endBlock == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	var err error
	fieldKeysType := [3]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2}
	for pos := 1; pos <= 3; pos++ {
		switch pos {
		case 1:
			if len(endBlock.ValidatorUpdates) == 0 {
				break
			}
			for i := 0; i < len(endBlock.ValidatorUpdates); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := endBlock.ValidatorUpdates[i].MarshalToAmino()
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		case 2:
			if endBlock.ConsensusParamUpdates == nil {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			data, err := MarshalConsensusParamsToAmino(*endBlock.ConsensusParamUpdates)
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSlice(&buf, data)
			if err != nil {
				return nil, err
			}
		case 3:
			eventsLen := len(endBlock.Events)
			if eventsLen == 0 {
				break
			}
			for i := 0; i < eventsLen; i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := endBlock.Events[i].MarshalToAmino()
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSlice(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}
