package quickproto

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), true)
	msg.Data = []byte("key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
	// Test and time the parsing of the message
	start_time := time.Now()
	_, err := msg.Parse()
	fmt.Println("(SHORT B64) Parse time:", time.Since(start_time))
	if err != nil {
		fmt.Println(err)
	}
	// Validate the message
	if msg.Headers["key1"][0] != "value1" {
		t.Error("(SHORT B64) Expected key1 to be value1 " + msg.Headers["key1"][0])
	}
	if msg.Headers["key1"][1] != "value2" {
		t.Error("(SHORT B64) Expected key1 to be value2 " + msg.Headers["key1"][1])
	}
	if msg.Headers["key2"][0] != "value2" {
		t.Error("(SHORT B64) Expected key2 to be value2 " + msg.Headers["key2"][0])
	}
	if string(msg.Body) != "BODYBODYBODY" {
		t.Error("(SHORT B64) Expected body to be BODYBODYBODY")
	}
}

func TestGenerate(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), true)
	msg.AddHeader("key1", "value1")
	msg.AddHeader("key1", "value2")
	msg.AddHeader("key2", "value2")
	msg.Body = []byte("BODYBODYBODY")
	// Test and time the parsing of the message
	start_time := time.Now()
	msg.Generate()
	fmt.Println("(SHORT B64) Generation time:", time.Since(start_time))
	// Validate the message
	if string(msg.Data) != "key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ" {
		if string(msg.Data) != "key2&value2&&key1&value1&value2&&&&Qk9EWUJPRFlCT0RZ" {
			t.Error("(SHORT B64)Expected data to be key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
			fmt.Println(string(msg.Data))
			fmt.Println("key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
			fmt.Println(string(msg.Data) == "key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
			fmt.Println(string(msg.Data) == "key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
			fmt.Println(len(string(msg.Data)), len("key1&value1&value2&&key2&value2&&&&Qk9EWUJPRFlCT0RZ"))
		}
	}
}

func TestParseLong(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), true)
	body := []byte(strings.Repeat("BODYBODYBODY_", 100000000)) // 13 * 100000000 = 1300000000 bytes (1300 MB)
	b64 := base64.StdEncoding.EncodeToString(body)
	msg.Data = []byte("key1&value1&&key2&value2&&&&" + b64)
	// Test and time the parsing of the message
	start_time := time.Now()
	_, err := msg.Parse()
	fmt.Println("(LONG B64) Parse time:", time.Since(start_time))
	if err != nil {
		t.Error(err)
	}
	// Validate the message
	if msg.Headers["key1"][0] != "value1" {
		t.Error("(LONG B64) Expected key1 to be value1 " + msg.Headers["key1"][0])
	}
	if msg.Headers["key2"][0] != "value2" {
		t.Error("(LONG B64) Expected key2 to be value2 " + msg.Headers["key2"][0])
	}
	if string(msg.Body) != strings.Repeat("BODYBODYBODY_", 100000000) { // 13 * 100000000 = 1300000000 bytes (1300 MB)
		t.Error("(LONG B64) Expected body to be BODYBODYBODY")
	}
}

func TestGenerateLong(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), true)
	msg.Headers["key1"] = []string{"value1"}
	msg.Headers["key2"] = []string{"value2"}
	// Create a 1300 MB body
	msg.Body = []byte(strings.Repeat("BODYBODYBODY_", 100000000)) // 13 * 100000000 = 1300000000 bytes (1300 MB)
	// Test and time the parsing of the message
	start_time := time.Now()
	msg.Generate()
	fmt.Println("(LONG B64) Generation time:", time.Since(start_time))
	// Validate the message
	b64 := base64.StdEncoding.EncodeToString(msg.Body)
	if string(msg.Data) != "key1&value1&&key2&value2&&&&"+b64 {
		if string(msg.Data) != "key2&value2&&key1&value1&&&&"+b64 {
			t.Error("(LONG B64) Expected data to be key1&value1&&key2&value2&&&&Qk9EWUJPRFlCT0RZ_")
		}
	}
}

func TestParse_NoB64(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), false)
	// msg.Data = []byte("key1&value1&&key2&value2&&&&Qk9EWUJPRFlCT0RZ")
	msg.Data = []byte("key1&value1&&key2&value2&&&&BODYBODYBODY")
	// Test and time the parsing of the message
	start_time := time.Now()
	_, err := msg.Parse()
	fmt.Println("(SHORT B64) Parse time:", time.Since(start_time))
	if err != nil {
		fmt.Println(err)
	}
	// Validate the message
	if msg.Headers["key1"][0] != "value1" {
		t.Error("(SHORT NO_B64) Expected key1 to be value1" + msg.Headers["key1"][0])
	}
	if msg.Headers["key2"][0] != "value2" {
		t.Error("(SHORT NO_B64) Expected key2 to be value2" + msg.Headers["key2"][0])
	}
	if string(msg.Body) != "BODYBODYBODY" {
		t.Error("(SHORT NO_B64) Expected body to be BODYBODYBODY")
	}
}

func TestGenerate_NoB64(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), false)
	msg.Headers["key1"] = []string{"value1"}
	msg.Headers["key2"] = []string{"value2"}
	msg.Body = []byte("BODYBODYBODY")
	// Test and time the parsing of the message
	start_time := time.Now()
	msg.Generate()
	fmt.Println("(SHORT NOB64) Generation time:", time.Since(start_time))
	// Validate the message
	if string(msg.Data) != "key1&value1&&key2&value2&&&&BODYBODYBODY" {
		if string(msg.Data) != "key2&value2&&key1&value1&&&&BODYBODYBODY" {
			t.Error("(SHORT NO_B64) Expected data to be key1&value1&&key2&value2&&&&BODYBODYBODY")
		}
	}
}

func TestParseLong_NoB64(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), false)
	body := []byte(strings.Repeat("BODYBODYBODY_", 100000000)) // 13 * 100000000 = 1300000000 bytes (1300 MB)
	// b64 := base64.StdEncoding.EncodeToString(body)
	msg.Data = []byte("key1&value1&&key2&value2&&&&" + string(body))
	// Test and time the parsing of the message
	start_time := time.Now()
	_, err := msg.Parse()
	fmt.Println("(LONG NO_B64) Parse time:", time.Since(start_time))
	if err != nil {
		t.Error(err)
	}
	// Validate the message
	if msg.Headers["key1"][0] != "value1" {
		t.Error("(LONG NO_B64) Expected key1 to be value1 " + msg.Headers["key1"][0])
	}
	if msg.Headers["key2"][0] != "value2" {
		t.Error("(LONG NO_B64) Expected key2 to be value2 " + msg.Headers["key2"][0])
	}
	if string(msg.Body) != strings.Repeat("BODYBODYBODY_", 100000000) { // 13 * 100000000 = 1300000000 bytes (1300 MB)
		t.Error("(LONG NO_B64) Expected body to be BODYBODYBODY")
	}
}

func TestGenerateLong_NoB64(t *testing.T) {
	// Create a new message
	msg := NewMessage([]byte("&"), false)
	msg.Headers["key1"] = []string{"value1"}
	msg.Headers["key2"] = []string{"value2"}
	// Create a 1300 MB body
	body := []byte(strings.Repeat("BODYBODYBODY_", 100000000)) // 13 * 100000000 = 1300000000 bytes (1300 MB)
	msg.Body = body
	// Test and time the parsing of the message
	start_time := time.Now()
	msg.Generate()
	fmt.Println("(LONG NO_B64) Generation time:", time.Since(start_time))
	// Validate the message
	// b64 := base64.StdEncoding.EncodeToString(msg.Body)
	if string(msg.Data) != "key1&value1&&key2&value2&&&&"+string(body) {
		if string(msg.Data) != "key2&value2&&key1&value1&&&&"+string(body) {
			t.Error("(LONG NO_B64) Expected data to be key1&value1&&key2&value2&&&&BODYBODYBODY_")
		}
	}
}

func TestGenerateAndParse(t *testing.T) {
	msg := NewMessage([]byte("###"), true)
	msg.AddHeader("key1", "value1")
	msg.AddHeader("key1", "value2")
	msg.AddHeader("key1", "value3")
	msg.AddHeader("key2", "value2")
	msg.Body = []byte("BODYBODYBODY")
	msg.Generate()
	_, err := msg.Parse()
	if err != nil {
		t.Error(err)
	}
	if msg.Headers["key1"][0] != "value1" {
		t.Error("Expected key1 to be value1")
	}
	if msg.Headers["key1"][1] != "value2" {
		t.Error("Expected key1 to be value2")
	}
	if msg.Headers["key1"][2] != "value3" {
		t.Error("Expected key1 to be value3")
	}
	if msg.Headers["key2"][0] != "value2" {
		t.Error("Expected key2 to be value2")
	}
	if string(msg.Body) != "BODYBODYBODY" {
		t.Error("Expected body to be BODYBODYBODY")
	}
}
