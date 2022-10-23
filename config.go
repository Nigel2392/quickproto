package quickproto

import (
	"bytes"
	"crypto/rsa"
)

type Config struct {
	Delimiter   []byte
	UseEncoding bool
	UseCrypto   bool
	BufSize     int
	Enc_func    func([]byte) []byte
	Dec_func    func([]byte) ([]byte, error)
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
}

// NewConfig creates a new Config.
func NewConfig(delimiter []byte, useencoding bool, usecrypto bool, bufsize int, enc_f func([]byte) []byte, dec_f func([]byte) ([]byte, error)) *Config {
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
		Enc_func:    enc_f,
		Dec_func:    dec_f,
	}
}

func (c *Config) NewMessage() *Message {
	return NewMessage(c.Delimiter, c.UseEncoding, c.Enc_func, c.Dec_func)
}
