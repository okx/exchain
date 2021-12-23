package tlv

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// golang encoding/binary official testcase of varint
var tests = []int64{
	-1 << 63,
	-1<<63 + 1,
	-1,
	0,
	1,
	2,
	10,
	20,
	63,
	64,
	65,
	127,
	128,
	129,
	255,
	256,
	257,
	1<<63 - 1,
}

// test the varint algo implementation
func TestVarintImpl(t *testing.T) {
	buf := make([]byte, 1024)
	buffer := NewBuffer()

	// Write layout test
	for _, v := range tests {
		buffer.encodeVarint(uint64(v))
		n := binary.PutUvarint(buf, uint64(v))
		assert.Equal(t, buf[:n], buffer.inner.Bytes())
		buffer.inner.Reset()
	}

	// WriteRead compare
	for _, v := range tests {
		buffer.encodeVarint(uint64(v))
	}
	for _, v := range tests {
		r := buffer.decodeVarint()
		assert.Equal(t, v, int64(r))
	}
}

func TestPrimitiveTypeReadWrite(t *testing.T) {
	buffer := NewBuffer()
	buffer.WriteInt16(int16(0x01))
	buffer.WriteInt32(int32(0x02))
	buffer.WriteInt64(int64(0x03))
	buffer.WriteUint16(uint16(0x04))
	buffer.WriteUint32(uint32(0x05))
	buffer.WriteUint64(uint64(0x06))
	buffer.WriteFloat32(float32(0.7))
	buffer.WriteFloat64(float64(0.8))

	v, ty := buffer.Read()
	assert.Equal(t, Int16, ty)
	assert.Equal(t, int16(0x01), v)

	v, ty = buffer.Read()
	assert.Equal(t, Int32, ty)
	assert.Equal(t, int32(0x02), v)

	v, ty = buffer.Read()
	assert.Equal(t, Int64, ty)
	assert.Equal(t, int64(0x03), v)

	v, ty = buffer.Read()
	assert.Equal(t, Uint16, ty)
	assert.Equal(t, uint16(0x04), v)

	v, ty = buffer.Read()
	assert.Equal(t, Uint32, ty)
	assert.Equal(t, uint32(0x05), v)

	v, ty = buffer.Read()
	assert.Equal(t, Uint64, ty)
	assert.Equal(t, uint64(0x06), v)

	v, ty = buffer.Read()
	assert.Equal(t, Float32, ty)
	assert.Equal(t, float32(0.7), v)

	v, ty = buffer.Read()
	assert.Equal(t, Float64, ty)
	assert.Equal(t, float64(0.8), v)

	_, ty = buffer.Read()
	assert.Equal(t, EmptyError, ty)
}

func TestByteArrayReadWrite(t *testing.T) {
	origin := []byte("test_slice///")
	l := len(origin)
	buffer := NewBuffer()
	buffer.Write(origin)

	assert.Equal(t, Bytes, Type(buffer.inner.Bytes()[0]))
	_, _ = buffer.inner.ReadByte()
	ol := buffer.decodeVarint()
	assert.Equal(t, l, int(ol))
	assert.Equal(t, origin, buffer.inner.Bytes())

	buffer.inner.Reset()
	buffer.Write(origin)
	v, ty := buffer.Read()
	assert.Equal(t, Bytes, ty)
	assert.Equal(t, origin, v)
}

func TestTlvErrorBranch(t *testing.T) {
	buffer := NewBuffer()
	buffer.WriteInt16(int16(0x01))
	origin := buffer.Bytes()
	origin[0] = 0xee
	buffer = With(origin)
	_, ty := buffer.Read()
	assert.Equal(t, UnsupportError, ty)
}
