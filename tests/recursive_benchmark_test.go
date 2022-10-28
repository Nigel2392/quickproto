package tests

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Nigel2392/quickproto"
)

func getGeneratedMessage(enc func(data []byte) []byte, dec func(data []byte) ([]byte, error)) *quickproto.Message {
	msg := quickproto.NewMessage([]byte("&"), true, enc, dec)
	msg.Headers = Preder_HEADERS_SMALL
	msg.Body = Predef_BODY_SMALL
	msg.Generate()
	return msg
}

func getPredefHeadersSmall() map[string][]string {
	curr_headers := make(map[string][]string)
	for i := 0; i < 50; i++ {
		curr_headers["key"+strconv.Itoa(i)] = append(curr_headers["key"+strconv.Itoa(i)], "value"+strconv.Itoa(i))
	}
	return curr_headers
}

var Preder_HEADERS_SMALL = getPredefHeadersSmall()
var Predef_BODY_SMALL = []byte(strings.Repeat("ABC", int(10000)))

const ITERS = 100000

func TestRecursiveGenerateB16(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base16Encoding, quickproto.Base16Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Generate()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestRecursiveGenerateB32(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base32Encoding, quickproto.Base32Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Generate()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestRecursiveGenerateB64(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base64Encoding, quickproto.Base64Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Generate()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestRecursiveGenerateGob(t *testing.T) {
	msg := getGeneratedMessage(quickproto.GobEncoding, quickproto.GobDecoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Generate()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestRecursiveGeneratePlain(t *testing.T) {
	msg := getGeneratedMessage(nil, nil)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Generate()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")
}

func TestRecursiveParseB16(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base16Encoding, quickproto.Base16Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Parse()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")

}

func TestRecursiveParseB32(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base32Encoding, quickproto.Base32Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Parse()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")

}

func TestRecursiveParseB64(t *testing.T) {
	msg := getGeneratedMessage(quickproto.Base64Encoding, quickproto.Base64Decoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Parse()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")

}

func TestRecursiveParseGob(t *testing.T) {
	msg := getGeneratedMessage(quickproto.GobEncoding, quickproto.GobDecoding)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Parse()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")

}

func TestRecursiveParsePlain(t *testing.T) {
	msg := getGeneratedMessage(nil, nil)
	start_time := time.Now()
	for i := 0; i < ITERS; i++ {
		msg.Parse()
	}
	t.Log(t.Name(), " [Size: "+strconv.Itoa(len(msg.Data))+" BYTES, Iterations:"+strconv.Itoa(ITERS)+"] finished in: ", time.Since(start_time).Milliseconds(), "ms")

}
