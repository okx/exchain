package types

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-amino"

	"github.com/pkg/errors"
)

// BlockMeta contains meta information.
type BlockMeta struct {
	BlockID   BlockID `json:"block_id"`
	BlockSize int     `json:"block_size"`
	Header    Header  `json:"header"`
	NumTxs    int     `json:"num_txs"`
}

// NewBlockMeta returns a new BlockMeta.
func NewBlockMeta(block *Block, blockParts *PartSet) *BlockMeta {
	return &BlockMeta{
		BlockID:   BlockID{block.Hash(), blockParts.Header()},
		BlockSize: block.Size(),
		Header:    block.Header,
		NumTxs:    len(block.Data.Txs),
	}
}

func (bm *BlockMeta) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
	const fieldCount = 4
	var currentField int
	var currentType amino.Typ3
	var err error

	for cur := 1; cur <= fieldCount; cur++ {
		if len(data) != 0 && (currentField == 0 || currentField < cur) {
			var nextField int
			if nextField, currentType, err = amino.ParseProtoPosAndTypeMustOneByte(data[0]); err != nil {
				return err
			}
			if nextField < currentField {
				return errors.Errorf("next field should greater than %d, got %d", currentField, nextField)
			} else {
				currentField = nextField
			}
		}

		if len(data) == 0 || currentField != cur {
			switch cur {
			case 1:
				bm.BlockID = BlockID{}
			case 2:
				bm.BlockSize = 0
			case 3:
				bm.Header = Header{}
			case 4:
				bm.NumTxs = 0
			default:
				return fmt.Errorf("unexpect feild num %d", cur)
			}
		} else {
			pbk := data[0]
			data = data[1:]
			var subData []byte
			if currentType == amino.Typ3_ByteLength {
				if subData, err = amino.DecodeByteSliceWithoutCopy(&data); err != nil {
					return err
				}
			}
			switch pbk {
			case 1<<3 | byte(amino.Typ3_ByteLength):
				if err = bm.BlockID.UnmarshalFromAmino(cdc, subData); err != nil {
					return err
				}
			case 2<<3 | byte(amino.Typ3_Varint):
				if bm.BlockSize, err = amino.DecodeIntUpdateBytes(&data); err != nil {
					return err
				}
			case 3<<3 | byte(amino.Typ3_ByteLength):
				if err = bm.Header.UnmarshalFromAmino(cdc, subData); err != nil {
					return err
				}
			case 4<<3 | byte(amino.Typ3_Varint):
				if bm.NumTxs, err = amino.DecodeIntUpdateBytes(&data); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpect pb key %d", pbk)
			}
		}
	}

	if len(data) != 0 {
		return errors.Errorf("unexpect data remain %X", data)
	}

	return nil
}

//-----------------------------------------------------------
// These methods are for Protobuf Compatibility

// Size returns the size of the amino encoding, in bytes.
func (bm *BlockMeta) Size() int {
	bs, _ := bm.Marshal()
	return len(bs)
}

// Marshal returns the amino encoding.
func (bm *BlockMeta) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(bm)
}

// MarshalTo calls Marshal and copies to the given buffer.
func (bm *BlockMeta) MarshalTo(data []byte) (int, error) {
	bs, err := bm.Marshal()
	if err != nil {
		return -1, err
	}
	return copy(data, bs), nil
}

// Unmarshal deserializes from amino encoded form.
func (bm *BlockMeta) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, bm)
}

// ValidateBasic performs basic validation.
func (bm *BlockMeta) ValidateBasic() error {
	if err := bm.BlockID.ValidateBasic(); err != nil {
		return err
	}
	if !bytes.Equal(bm.BlockID.Hash, bm.Header.Hash()) {
		return errors.Errorf("expected BlockID#Hash and Header#Hash to be the same, got %X != %X",
			bm.BlockID.Hash, bm.Header.Hash())
	}
	return nil
}
