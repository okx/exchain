package types

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/kv"

	"github.com/tendermint/go-amino"
)

func (pubkey PubKey) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
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

// UnmarshalFromAmino unmarshal data from amino bytes.
func (pub *PubKey) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			pub.Type = string(subData)

		case 2:
			pub.Data = make([]byte, len(subData))
			copy(pub.Data, subData)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (valUpdate ValidatorUpdate) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [2]byte{1<<3 | 2, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			var data []byte
			data, err = valUpdate.PubKey.MarshalToAmino(cdc)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
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
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from amino bytes.
func (vu *ValidatorUpdate) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			err := vu.PubKey.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}

		case 2:
			power, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			vu.Power = int64(power)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (params BlockParams) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1 << 3, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		var err error
		switch pos {
		case 1:
			if params.MaxBytes == 0 {
				break
			}
			err = amino.EncodeUvarintWithKeyToBuffer(&buf, uint64(params.MaxBytes), fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
		case 2:
			if params.MaxGas == 0 {
				break
			}
			err = amino.EncodeUvarintWithKeyToBuffer(&buf, uint64(params.MaxGas), fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from amino bytes.
func (bp *BlockParams) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	var dataLen uint64 = 0

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, _, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		switch pos {
		case 1:
			mb, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			bp.MaxBytes = int64(mb)
		case 2:
			mg, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			bp.MaxGas = int64(mg)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (params EvidenceParams) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1 << 3, 2 << 3}
	for pos := 1; pos <= 2; pos++ {
		var err error
		switch pos {
		case 1:
			if params.MaxAgeNumBlocks == 0 {
				break
			}
			err = amino.EncodeUvarintWithKeyToBuffer(&buf, uint64(params.MaxAgeNumBlocks), fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
		case 2:
			if params.MaxAgeDuration == 0 {
				break
			}
			err = amino.EncodeUvarintWithKeyToBuffer(&buf, uint64(params.MaxAgeDuration), fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from amino bytes.
func (ep *EvidenceParams) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	var dataLen uint64 = 0

	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}
		pos, _, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		switch pos {
		case 1:
			ma, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			ep.MaxAgeNumBlocks = int64(ma)

		case 2:
			md, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			ep.MaxAgeDuration = time.Duration(md)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (params ValidatorParams) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	var pubKeyTypesPbKey = byte(1<<3 | 2)
	for i := 0; i < len(params.PubKeyTypes); i++ {
		err = buf.WriteByte(pubKeyTypesPbKey)
		if err != nil {
			return nil, err
		}
		err = amino.EncodeStringToBuffer(&buf, params.PubKeyTypes[i])
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from amino bytes.
func (vp *ValidatorParams) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			vp.PubKeyTypes = append(vp.PubKeyTypes, string(subData))

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (event Event) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if event.Type == "" {
				break
			}
			err = amino.EncodeStringWithKeyToBuffer(&buf, event.Type, fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
		case 2:
			for i := 0; i < len(event.Attributes); i++ {
				data, err := event.Attributes[i].MarshalToAmino(cdc)
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceWithKeyToBuffer(&buf, data, fieldKeysType[pos-1])
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

func (event *Event) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
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
			err = kvpair.UnmarshalFromAmino(cdc, subData)
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

func (tx *ResponseDeliverTx) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
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
				data, err := tx.Events[i].MarshalToAmino(cdc)
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

func (tx *ResponseDeliverTx) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
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
			dataLen, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}

			data = data[n:]
			if len(data) < int(dataLen) {
				return errors.New("invalid datalen")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
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
			if err != nil {
				return err
			}
			tx.GasWanted = int64(uvint)
			dataLen = uint64(n)
		case 6:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			tx.GasUsed = int64(uvint)
			dataLen = uint64(n)
		case 7:
			var event Event
			err = event.UnmarshalFromAmino(cdc, subData)
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

func (beginBlock ResponseBeginBlock) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKey := byte(1<<3 | 2)
	for i := 0; i < len(beginBlock.Events); i++ {
		err := buf.WriteByte(fieldKey)
		if err != nil {
			return nil, err
		}
		data, err := beginBlock.Events[i].MarshalToAmino(cdc)
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

func (bb *ResponseBeginBlock) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			event := Event{}
			err = event.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			bb.Events = append(bb.Events, event)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (params ConsensusParams) MarshalToAmino(cdc *amino.Codec) (data []byte, err error) {
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
			data, err = params.Block.MarshalToAmino(cdc)
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
			data, err = params.Evidence.MarshalToAmino(cdc)
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
			data, err = params.Validator.MarshalToAmino(cdc)
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
func (cp *ConsensusParams) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			bParams := &BlockParams{}
			err := bParams.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			cp.Block = bParams
		case 2:
			eParams := &EvidenceParams{}
			err := eParams.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			cp.Evidence = eParams
		case 3:
			vp := &ValidatorParams{}
			err := vp.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			cp.Validator = vp

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (endBlock ResponseEndBlock) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	fieldKeysType := [3]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2}
	for pos := 1; pos <= 3; pos++ {
		switch pos {
		case 1:
			for i := 0; i < len(endBlock.ValidatorUpdates); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := endBlock.ValidatorUpdates[i].MarshalToAmino(cdc)
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
			data, err := endBlock.ConsensusParamUpdates.MarshalToAmino(cdc)
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSlice(&buf, data)
			if err != nil {
				return nil, err
			}
		case 3:
			for i := 0; i < len(endBlock.Events); i++ {
				err = buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := endBlock.Events[i].MarshalToAmino(cdc)
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

func (eb *ResponseEndBlock) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			vu := ValidatorUpdate{}
			err := vu.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			eb.ValidatorUpdates = append(eb.ValidatorUpdates, vu)
		case 2:
			consParam := &ConsensusParams{}
			err := consParam.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			eb.ConsensusParamUpdates = consParam
		case 3:
			var event Event
			err = event.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			eb.Events = append(eb.Events, event)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
