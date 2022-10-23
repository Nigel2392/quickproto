package quickproto

import (
	"bytes"
	"crypto/rsa"
)

// General configuration to use for client and server.
type Config struct {
	// Delimiter used for separating message data.
	Delimiter []byte
	// Use encoding?
	UseEncoding bool
	// Use crypto?
	UseCrypto bool
	// Buffer size.
	BufSize int
	// Encoding/Decoding functions.
	Encode_func func([]byte) []byte
	Decode_func func([]byte) ([]byte, error)
	// RSA keys
	PrivateKey *rsa.PrivateKey // Server-side.
	PublicKey  *rsa.PublicKey  // Client-side.
}

// NewConfig creates a new Config.
func NewConfig(delimiter []byte, useencoding bool, usecrypto bool, bufsize int, encode_f func([]byte) []byte, decode_f func([]byte) ([]byte, error)) *Config {
	if delimiter == nil {
		delimiter = STANDARD_DELIM
	}
	for _, d := range BANNED_DELIMITERS {
		if bytes.Contains(delimiter, []byte(d)) {
			panic("Delimiter contains banned characters: " + d)
		}
	}
	return &Config{
		Delimiter:   delimiter,
		UseEncoding: useencoding,
		UseCrypto:   usecrypto,
		BufSize:     bufsize,
		Encode_func: encode_f,
		Decode_func: decode_f,
	}
}

// Generate a new message with default configuration options.
func (c *Config) NewMessage() *Message {
	return NewMessage(c.Delimiter, c.UseEncoding, c.Encode_func, c.Decode_func)
}
