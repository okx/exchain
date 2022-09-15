package store

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-amino"
)

func unmarshalBlockPartBytesTo(data []byte, buf *bytes.Buffer) error {
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
				return fmt.Errorf("not enough data for %s, need %d, have %d", aminoType, dataLen, len(data))
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			_, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
		case 2:
			buf.Write(subData)
			return nil
		case 3:
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
