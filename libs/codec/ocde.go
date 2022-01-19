package ocdc

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"sync"
)
type OCDC_TYPE int

const (
	RLP  OCDC_TYPE = iota
	JSON
	AMINO
)

var (
	once sync.Once
	ocdcType OCDC_TYPE = 0
)

func InitOcdc(cdcType OCDC_TYPE)  {
	once.Do(func() {
		ocdcType = cdcType
	})
}

func Encode(val interface{}) ([]byte, error) {
	switch ocdcType {
	case RLP:
		return rlp.EncodeToBytes(val)
	case JSON:
		return json.Marshal(val)
	}

	return nil, fmt.Errorf("unknown ocdc type")
}

func Decode(b []byte, val interface{}) error {
	switch ocdcType {
	case RLP:
		return rlp.DecodeBytes(b, val)
	case JSON:
		return json.Unmarshal(b, val)

	}
	return fmt.Errorf("unknown ocdc type")
}

