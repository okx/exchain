package go_benchmark

import (
	"encoding/json"
	"fmt"
	"testing"

	gojson "github.com/goccy/go-json"
	"github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

// compare output between std json and gojson
func TestJsonOutput(t *testing.T) {
	var d1, d2 SmallPayload

	gojson.Unmarshal(smallFixture, &d1)
	fmt.Println("gojson unm:\n", d1)

	json.Unmarshal(smallFixture, &d2)
	fmt.Println("std unm:\n", d2)

	v := SmallPayload{
		11, 22, "33", 44, "abcd", "abcd2", "abcd3", 55, 66,
	}

	t1, _ := gojson.Marshal(v)
	fmt.Println("gojson marshal:\n", string(t1))

	t2, _ := json.Marshal(v)
	fmt.Println("json marshal:\n", string(t2))

}

// fastjson
func BenchmarkFastJsonUnmarshal(b *testing.B) {
	b.Run("fastjson-small", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var p fastjson.Parser
			var data SmallPayload
			v, _ := p.ParseBytes(smallFixture)
			data.Uuid = string(v.GetStringBytes("uuid"))
			data.Tz = v.GetInt("tz")
			data.Ua = string(v.GetStringBytes("ua"))
			data.St = v.GetInt("st")
		}
	})

	b.Run("fastjson-medium", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var p fastjson.Parser
			var data MediumPayload
			v, _ := p.ParseBytes(mediumFixture)
			v.GetObject("person")
			data.Company = string(v.GetStringBytes("company"))
		}
	})

	b.Run("fastjson-large", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(largeFixture)))
		for i := 0; i < b.N; i++ {
			var p fastjson.Parser
			var data LargePayload
			v, _ := p.ParseBytes(largeFixture)
			v.GetObject("users")
			v.GetObject("topcs")
			_ = data
		}
	})
}

//jsoniter
func BenchmarkJsoniterUnmarshal(b *testing.B) {
	b.Run("jsoniter-small", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data SmallPayload
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			json.Unmarshal(smallFixture, &data)
		}
	})

	b.Run("jsoniter-medium", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data MediumPayload
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			json.Unmarshal(mediumFixture, &data)
		}
	})

	b.Run("jsoniter-large", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(largeFixture)))
		for i := 0; i < b.N; i++ {
			var data LargePayload
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			json.Unmarshal(largeFixture, &data)
		}
	})
}

// go-json
func BenchmarkGoJsonUnmarshal(b *testing.B) {
	b.Run("gojson-small", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data SmallPayload
			gojson.Unmarshal(smallFixture, &data)
		}
	})

	b.Run("gojson-medium", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data MediumPayload
			gojson.Unmarshal(mediumFixture, &data)
		}
	})

	b.Run("gojson-large", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(largeFixture)))
		for i := 0; i < b.N; i++ {
			var data LargePayload
			gojson.Unmarshal(largeFixture, &data)
		}
	})
}

/*
   std json
*/
func BenchmarkSTDJsonUnmarshal(b *testing.B) {
	b.Run("stdjson-small", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data SmallPayload
			json.Unmarshal(smallFixture, &data)
		}
	})

	b.Run("stdjson-medium", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var data MediumPayload
			json.Unmarshal(mediumFixture, &data)
		}
	})

	b.Run("stdjson-large", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(largeFixture)))
		for i := 0; i < b.N; i++ {
			var data LargePayload
			json.Unmarshal(largeFixture, &data)
		}
	})
}
