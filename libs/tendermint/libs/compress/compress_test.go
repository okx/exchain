package compress

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {

	compressType := 3
	compressFlag := 2

	data := []byte("test compress 2021")

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

func TestCompressTable(t *testing.T) {
	const compressMethod = 4
	type resultMethod [compressMethod]int
	const (
		zlib  = 1
		flate = 2
		gzip  = 3
		dummy = 4

		fastMode = 1
		bestMode = 2
		defaMode = 3

		actSuc = 0
		actErr = 1
		actRec = 2
		actDum = 3
	)

	testCompressShort := []byte("test compress short")
	testCompressLong := []byte("this is a long long long byte array. it's really very very long. we use the long byte array to test different compress method, the behavior should be the save with the short one above.")
	type args struct {
		compressType   int
		flag           int
		src            []byte
		unCompressType []int
	}
	tests := []struct {
		name string
		args args
		want []byte
		ret  resultMethod
	}{
		{"zlib      fast", args{compressType: zlib, flag: fastMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actSuc, actErr, actRec, actDum}},
		{"zlib      best", args{compressType: zlib, flag: bestMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actSuc, actErr, actRec, actDum}},
		{"zlib   default", args{compressType: zlib, flag: defaMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actSuc, actErr, actRec, actDum}},
		{"flate    fast", args{compressType: flate, flag: fastMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actSuc, actRec, actDum}},
		{"flate    best", args{compressType: flate, flag: bestMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actSuc, actRec, actDum}},
		{"falte default", args{compressType: flate, flag: defaMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actSuc, actRec, actDum}},
		{"gzip      fast", args{compressType: gzip, flag: fastMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actSuc, actDum}},
		{"gzip      best", args{compressType: gzip, flag: bestMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actSuc, actDum}},
		{"gzip   default", args{compressType: gzip, flag: defaMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actSuc, actDum}},
		{"dummy    fast", args{compressType: dummy, flag: fastMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actRec, actSuc}},
		{"dummy    best", args{compressType: dummy, flag: bestMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actRec, actDum}},
		{"dummy default", args{compressType: dummy, flag: defaMode, src: testCompressShort, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressShort, resultMethod{actRec, actErr, actRec, actDum}},

		{"zlib      fast long", args{compressType: zlib, flag: fastMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actSuc, actErr, actRec, actDum}},
		{"zlib      best long", args{compressType: zlib, flag: bestMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actSuc, actErr, actRec, actDum}},
		{"zlib   default long", args{compressType: zlib, flag: defaMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actSuc, actErr, actRec, actDum}},
		{"flate    fast long", args{compressType: flate, flag: fastMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actSuc, actRec, actDum}},
		{"flate    best long", args{compressType: flate, flag: bestMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actSuc, actRec, actDum}},
		{"falte default long", args{compressType: flate, flag: defaMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actSuc, actRec, actDum}},
		{"gzip      fast long", args{compressType: gzip, flag: fastMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actSuc, actDum}},
		{"gzip      best long", args{compressType: gzip, flag: bestMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actSuc, actDum}},
		{"gzip   default long", args{compressType: gzip, flag: defaMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actSuc, actDum}},
		{"dummy    fast long", args{compressType: dummy, flag: fastMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actRec, actSuc}},
		{"dummy    best long", args{compressType: dummy, flag: bestMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actRec, actDum}},
		{"dummy default long", args{compressType: dummy, flag: defaMode, src: testCompressLong, unCompressType: []int{zlib, flate, gzip, dummy}}, testCompressLong, resultMethod{actRec, actErr, actRec, actDum}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.args.compressType, tt.args.flag, tt.args.src)
			assert.NoError(t, err)
			for i, v := range tt.ret {
				switch v {
				case actSuc:
					restored, err := UnCompress(tt.args.unCompressType[i], got)
					assert.Nil(t, err)
					assert.Equalf(t, tt.want, restored, "Compress(%v, %v, %v)", tt.args.compressType, tt.args.flag, tt.args.src)
				case actErr:
					_, err := UnCompress(tt.args.unCompressType[i], got)
					assert.Error(t, err)
				case actRec:
					_, err := UnCompress(tt.args.unCompressType[i], got)
					assert.ErrorContainsf(t, err, "uncompress panic", err.Error())
				case actDum:
					restored, err := UnCompress(tt.args.unCompressType[i], got)
					assert.NoError(t, err)
					assert.Equalf(t, got, restored, "dummy compress(%v, %v, %v)", tt.args.compressType, tt.args.flag, tt.args.src)
				}
			}
		})
	}
}
