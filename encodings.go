package quickproto

import (
	"bytes"
	"encoding/base32"
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

// Base32 encoding
func Base32Encoding(data []byte) []byte {
	var b32_buffer bytes.Buffer
	encoder := base32.NewEncoder(base32.StdEncoding, &b32_buffer)
	encoder.Write(data)
	encoder.Close()
	return b32_buffer.Bytes()
}

// Base32 decoding
func Base32Decoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := base32.NewDecoder(base32.StdEncoding, buf)
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
