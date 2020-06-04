package keeper

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncompress(t *testing.T) {
	wasmRaw, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	wasmGzipped, err := ioutil.ReadFile("./testdata/contract.wasm.gzip")
	require.NoError(t, err)

	specs := map[string]struct {
		src       []byte
		expError  error
		expResult []byte
	}{
		"handle wasm uncompressed": {
			src:       wasmRaw,
			expResult: wasmRaw,
		},
		"handle wasm compressed": {
			src:       wasmGzipped,
			expResult: wasmRaw,
		},
		"handle nil slice": {
			src:       nil,
			expResult: nil,
		},
		"handle short unidentified": {
			src:       []byte{0x1, 0x2},
			expResult: []byte{0x1, 0x2},
		},
		"handle big input slice": {
			src:       []byte(strings.Repeat("a", maxSize+1)),
			expResult: []byte(strings.Repeat("a", maxSize+1)),
		},
		"handle gzip identifier only": {
			src:      gzipIdent,
			expError: io.ErrUnexpectedEOF,
		},
		"handle broken gzip": {
			src:      append(gzipIdent, byte(0x1)),
			expError: io.ErrUnexpectedEOF,
		},
		"handle incomplete gzip": {
			src:      wasmGzipped[:len(wasmGzipped)-5],
			expError: io.ErrUnexpectedEOF,
		},
		"handle big gzip output": {
			src:      asGzip(strings.Repeat("a", maxSize+1)),
			expError: io.ErrUnexpectedEOF,
		},
		"handle other big gzip output": {
			src:      asGzip(strings.Repeat("a", 2*maxSize)),
			expError: io.ErrUnexpectedEOF,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			r, err := uncompress(spec.src)
			if spec.expError != nil {
				return
			}

			require.True(t, err == spec.expError, "exp %+v got %+v", spec.expError, err)
			assert.Equal(t, spec.expResult, r)
		})
	}

}

func asGzip(src string) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(gzip.NewWriter(&buf), strings.NewReader(src)); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
