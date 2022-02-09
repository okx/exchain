package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"sync"
)

type compressBroker interface {
	defaultCompress(src []byte) ([]byte, error)
	bestCompress(src []byte) ([]byte, error)
	fastCompress(src []byte) ([]byte, error)
	unCompress(compressSrc []byte) ([]byte, error)
}

var (
	once                sync.Once
	zlibCompressBroker  compressBroker
	flateCompressBroker compressBroker
	gzipCompressBroker  compressBroker
	dummyCompressBroker compressBroker
)

func init() {
	once.Do(func() {
		zlibCompressBroker = &Zlib{}
		flateCompressBroker = &Flate{}
		gzipCompressBroker = &Gzip{}
		dummyCompressBroker = &dummy{}
	})
}

func Compress(compressType, flag int, src []byte) ([]byte, error) {
	bk := getCompressBroker(compressType)
	switch flag {
	case 1:
		return bk.fastCompress(src)
	case 2:
		return bk.bestCompress(src)
	default:
	}
	return bk.defaultCompress(src)
}

func UnCompress(compressType int, src []byte) ([]byte, error) {
	bk := getCompressBroker(compressType)
	return bk.unCompress(src)
}

func getCompressBroker(compressType int) compressBroker {
	var broker compressBroker
	switch compressType {
	case 1:
		broker = zlibCompressBroker
		break
	case 2:
		broker = flateCompressBroker
		break
	case 3:
		broker = gzipCompressBroker
		break
	default:
		broker = dummyCompressBroker
	}
	return broker
}

type dummy struct {
}

func (z *dummy) defaultCompress(src []byte) ([]byte, error) { return src, nil }
func (z *dummy) bestCompress(src []byte) ([]byte, error)    { return src, nil }
func (z *dummy) fastCompress(src []byte) ([]byte, error)    { return src, nil }
func (z *dummy) unCompress(src []byte) ([]byte, error)      { return src, nil }

// Zlib
type Zlib struct {
}

func (z *Zlib) defaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) bestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) fastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (z *Zlib) unCompress(compressSrc []byte) ([]byte, error) {
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

func (g *Gzip) defaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) bestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) fastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (g *Gzip) unCompress(compressSrc []byte) ([]byte, error) {
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

func (f *Flate) defaultCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) bestCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestCompression)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) fastCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestSpeed)
	w.Write(src)
	w.Close()
	return in.Bytes(), err
}

func (f *Flate) unCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := flate.NewReader(b)
	_, err := io.Copy(&out, r)
	return out.Bytes(), err
}
