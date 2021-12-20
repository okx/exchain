package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

type CompressBroker interface {
	DefaultCompress(src []byte) ([]byte, error)
	BestCompress(src []byte) ([]byte, error)
	FastCompress(src []byte) ([]byte, error)
	UnCompress(compressSrc []byte) ([]byte, error)
}

// Zlib
type Zlib struct {
}

func (z *Zlib) DefaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) BestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) FastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) UnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes(), err
}

// -------------------------------------------------------------

// Gzip
type Gzip struct {
}

func (g *Gzip) DefaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) BestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) FastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) UnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := gzip.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes(), err
}

// -------------------------------------------------------------

// Flate
type Flate struct {
}

func (f *Flate) DefaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) BestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) FastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) UnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := flate.NewReader(b)
	_, err := io.Copy(&out, r)
	return out.Bytes(), err
}