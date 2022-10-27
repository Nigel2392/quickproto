//go:build !benchmark
// +build !benchmark

package tests

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Nigel2392/quickproto"
)

func getGenerated(enc func(data []byte) []byte, dec func(data []byte) ([]byte, error)) *quickproto.Message {
	msg := quickproto.NewMessage([]byte("&"), true, enc, dec)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	msg.Generate()
	return msg
}

func getPredefHeaders() map[string][]string {
	curr_headers := make(map[string][]string)
	for i := 0; i < 1000000; i++ {
		curr_headers["key"+strconv.Itoa(i)] = append(curr_headers["key"+strconv.Itoa(i)], "value"+strconv.Itoa(i))
	}
	return curr_headers
}

var Predef_HEADERS = getPredefHeaders()
var Predef_BODY = []byte(strings.Repeat("ABC", int(1000000000/3)))

func TestB64Generate(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base64Encoding, quickproto.Base64Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	start_time := time.Now()
	msg.Generate()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestB32Generate(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base32Encoding, quickproto.Base32Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	start_time := time.Now()
	msg.Generate()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestHexGenerate(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base16Encoding, quickproto.Base16Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	start_time := time.Now()
	msg.Generate()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestGobGenerate(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.GobEncoding, quickproto.GobDecoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	start_time := time.Now()
	msg.Generate()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestPlainGenerate(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), false, nil, nil)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	start_time := time.Now()
	msg.Generate()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestB64Parse(t *testing.T) {
	msg := getGenerated(quickproto.Base64Encoding, quickproto.Base64Decoding)
	start_time := time.Now()
	msg.Parse()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestB32Parse(t *testing.T) {
	msg := getGenerated(quickproto.Base32Encoding, quickproto.Base32Decoding)
	start_time := time.Now()
	msg.Parse()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB]finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestHexParse(t *testing.T) {
	msg := getGenerated(quickproto.Base16Encoding, quickproto.Base16Decoding)
	start_time := time.Now()
	msg.Parse()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB]finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestGobParse(t *testing.T) {
	msg := getGenerated(quickproto.GobEncoding, quickproto.GobDecoding)
	start_time := time.Now()
	msg.Parse()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB]finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestPlainParse(t *testing.T) {
	msg := quickproto.NewMessage([]byte("&"), false, nil, nil)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	msg.Generate()
	start_time := time.Now()
	msg.Parse()
	t.Log(t.Name(), "[Size: "+strconv.Itoa(len(msg.Data)/1024/1024)+"MB]finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func BenchmarkB64Generate(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base64Encoding, quickproto.Base64Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Generate()
	}
}

func BenchmarkB32Generate(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base32Encoding, quickproto.Base32Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Generate()
	}
}

func BenchmarkHexGenerate(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.Base16Encoding, quickproto.Base16Decoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Generate()
	}
}

func BenchmarkGobGenerate(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), true, quickproto.GobEncoding, quickproto.GobDecoding)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Generate()
	}
}

func BenchmarkPlainGenerate(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), false, nil, nil)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Generate()
	}
}

func BenchmarkB64Parse(b *testing.B) {
	msg := getGenerated(quickproto.Base64Encoding, quickproto.Base64Decoding)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Parse()
	}
}

func BenchmarkB32Parse(b *testing.B) {
	msg := getGenerated(quickproto.Base32Encoding, quickproto.Base32Decoding)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Parse()
	}
}

func BenchmarkHexParse(b *testing.B) {
	msg := getGenerated(quickproto.Base16Encoding, quickproto.Base16Decoding)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Parse()
	}
}

func BenchmarkGobParse(b *testing.B) {
	msg := getGenerated(quickproto.GobEncoding, quickproto.GobDecoding)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Parse()
	}
}

func BenchmarkPlainParse(b *testing.B) {
	msg := quickproto.NewMessage([]byte("&"), false, nil, nil)
	msg.Headers = Predef_HEADERS
	msg.Body = Predef_BODY
	msg.Generate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg.Parse()
	}
}
