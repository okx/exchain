package compress

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {

	var tools = []CompressBroker{
		&Zlib{},
		&Gzip{},
		&Flate{},
	}

	data := []byte("oec 2021")

	for _, tool := range tools {
		res, err := tool.DefaultCompress(data)
		assert.Nil(t, err)
		verify := func(compressed []byte) {
			unCompressresed, err := tool.UnCompress(compressed)
			assert.Nil(t, err)
			assert.Equal(t, 0, bytes.Compare(data, unCompressresed))
		}
		verify(res)

		res, err = tool.BestCompress(data)
		assert.Nil(t, err)
		verify(res)

		res, err = tool.FastCompress(data)
		assert.Nil(t, err)
		verify(res)
	}
}
