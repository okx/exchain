package compress

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {

	compressType := 3
	compressFlag := 2

	data := []byte("okc 2021")

	for ctype := 0; ctype <= compressType; ctype++ {
		for flag := 0; flag <= compressFlag; flag++ {
			res, err := Compress(ctype, flag, data)
			assert.Nil(t, err)
			unCompressresed, err := UnCompress(ctype, res)
			assert.Nil(t, err)
			assert.Equal(t, 0, bytes.Compare(data, unCompressresed))
		}
	}
}
