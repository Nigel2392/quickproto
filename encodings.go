package quickproto

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// Base64 encoding
func Base64Encoding(data []byte) []byte {
	var b64_buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b64_buffer)
	encoder.Write(data)
	encoder.Close()
	return b64_buffer.Bytes()
}

// Base64 decoding
func Base64Decoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	return io.ReadAll(decoder)
}

// Hex encoding
func Base16Encoding(data []byte) []byte {
	return []byte(hex.EncodeToString(data))
}

// Hex decoding
func Base16Decoding(data []byte) ([]byte, error) {
	return hex.DecodeString(string(data))
}
