package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

type CompressBroker interface {
	DefaultCompress(src []byte) []byte
	BestCompress(src []byte) []byte
	FastCompress(src []byte) []byte
	UnCompress(compressSrc []byte) []byte
}

// Zlib
type Zlib struct {
}

func (z *Zlib) DefaultCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (z *Zlib) BestCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (z *Zlib) FastCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (z *Zlib) UnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

// -------------------------------------------------------------

// Gzip
type Gzip struct {
}

func (g *Gzip) DefaultCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (g *Gzip) BestCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (g *Gzip) FastCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (g *Gzip) UnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := gzip.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

// -------------------------------------------------------------

// Flate
type Flate struct {
}

func (f *Flate) DefaultCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := flate.NewWriter(&in, flate.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (f *Flate) BestCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := flate.NewWriter(&in, flate.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (f *Flate) FastCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := flate.NewWriter(&in, flate.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func (f *Flate) UnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := flate.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}