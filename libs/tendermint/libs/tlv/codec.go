package tlv

import (
	"bytes"
	"encoding/binary"
	"sync"
)

// tlv is a simple lib used to generate TLV formate codec for any structure
// you have to defined your own logic handling your own structure
// NOTE: you have to keep your write buffer order and read order from buffer

// objects pool to optimize the GC
var pool sync.Pool

func init() {
	pool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

}

// Type represents the underlying data description
type Type byte

const (
	Int16   Type = 0x01
	Int32   Type = 0x02
	Int64   Type = 0x03
	Uint16  Type = 0x04
	Uint32  Type = 0x05
	Uint64  Type = 0x06
	Float32 Type = 0x07
	Float64 Type = 0x08
	Bytes   Type = 0x09

	// errors of this buffer
	None           Type = 0xf0
	EmptyError     Type = 0xf1
	UnsupportError Type = 0xf2
	ParseError     Type = 0xf3
)

// Buffer to implement the self-defined TLV codec
type Buffer struct {
	inner *bytes.Buffer
}

// With current buffer to create a new Buffer for reading
func With(buffer []byte) *Buffer {
	return &Buffer{
		inner: bytes.NewBuffer(buffer),
	}
}

// NewBuffer create a new Buffer for writing
func NewBuffer() *Buffer {
	v := pool.Get()
	return &Buffer{
		inner: v.(*bytes.Buffer),
	}
}

// Bytes serialize curretn Buffer to []byte
// after this method the underlying buffer will be reset
// you should use this method as final method you call
func (buf *Buffer) Bytes() []byte {
	defer func() {
		// to keep the underlying buffer empty semantic correctly
		buf.inner.Reset()
		pool.Put(buf.inner)
	}()

	return buf.inner.Bytes()
}

// encode v into varint and write into inner directly
func (buffer *Buffer) encodeVarint(v uint64) {
	for n := 0; v > 127; n++ {
		buffer.inner.WriteByte(0x80 | uint8(v&0x7F))
		v >>= 7
	}
	buffer.inner.WriteByte(uint8(v))
	return
}

// decode varint from the inner and return v
func (buffer *Buffer) decodeVarint() (v uint64) {
	for shift := uint(0); shift < 64; shift += 7 {
		if buffer.inner.Len() == 0 {
			return
		}
		b, _ := buffer.inner.ReadByte() // MUST success if not should panic right now
		v |= (uint64(b) & 0x7F) << shift
		if (uint64(b) & 0x80) == 0 {
			return
		}
	}
	return
}

//==================setters==================

// nolint
func (buf *Buffer) WriteInt16(v int16) *Buffer {
	buf.inner.WriteByte(byte(Int16))
	buf.encodeVarint(uint64(v))
	return buf
}

// nolint
func (buf *Buffer) WriteInt32(v int32) *Buffer {
	buf.inner.WriteByte(byte(Int32))
	buf.encodeVarint(uint64(v))
	return buf
}

// nolint
func (buf *Buffer) WriteInt64(v int64) *Buffer {
	buf.inner.WriteByte(byte(Int64))
	buf.encodeVarint(uint64(v))
	return buf
}

// nolint
func (buf *Buffer) WriteUint16(v uint16) *Buffer {
	buf.inner.WriteByte(byte(Uint16))
	buf.encodeVarint(uint64(v))
	return buf
}

// nolint
func (buf *Buffer) WriteUint32(v uint32) *Buffer {
	buf.inner.WriteByte(byte(Uint32))
	buf.encodeVarint(uint64(v))
	return buf
}

// nolint
func (buf *Buffer) WriteUint64(v uint64) *Buffer {
	buf.inner.WriteByte(byte(Uint64))
	buf.encodeVarint(v)
	return buf
}

// nolint
func (buf *Buffer) WriteFloat32(v float32) *Buffer {
	buf.inner.WriteByte(byte(Float32))
	binary.Write(buf.inner, binary.BigEndian, v)
	return buf
}

// nolint
func (buf *Buffer) WriteFloat64(v float64) *Buffer {
	buf.inner.WriteByte(byte(Float64))
	binary.Write(buf.inner, binary.BigEndian, v)
	return buf
}

// nolint
func (buf *Buffer) Write(v []byte) *Buffer {
	buf.inner.WriteByte(byte(Bytes))
	buf.encodeVarint(uint64(len(v)))
	buf.inner.Write(v)
	return buf
}

//==================getters==================

// Read is the generic method to read the data from underlying buffer
// if the underlying buffer readable is 0 will return EmptyError
// if the Type of current type is unsupported return UnsupportError
func (buf *Buffer) Read() (interface{}, Type) {
	if buf.inner.Len() == 0 {
		return nil, EmptyError
	}
	ty, _ := buf.inner.ReadByte()
	switch Type(ty) {
	case Int16:
		return int16(buf.decodeVarint()), Int16
	case Int32:
		return int32(buf.decodeVarint()), Int32
	case Int64:
		return int64(buf.decodeVarint()), Int64
	case Uint16:
		return uint16(buf.decodeVarint()), Uint16
	case Uint32:
		return uint32(buf.decodeVarint()), Uint32
	case Uint64:
		return buf.decodeVarint(), Uint64
	case Float32:
		{
			var f float32
			e := binary.Read(buf.inner, binary.BigEndian, &f)
			if e != nil {
				f = 0.0
				return f, EmptyError
			}
			return f, Float32
		}
	case Float64:
		{
			var f float64
			e := binary.Read(buf.inner, binary.BigEndian, &f)
			if e != nil {
				f = 0.0
				return f, EmptyError
			}
			return f, Float64
		}
	case Bytes:
		{
			l := buf.decodeVarint()
			b := make([]byte, int(l))
			buf.inner.Read(b)
			return b, Bytes
		}
	default:
		return nil, UnsupportError
	}
}
