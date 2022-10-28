package quickproto

import (
	"bytes"
	"compress/gzip"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"io"
)

// Base 64 encoding.
func Base64Encoding(data []byte) []byte {
	var b64_buffer bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b64_buffer)
	encoder.Write(data)
	encoder.Close()
	return b64_buffer.Bytes()
}

// Base 64 decoding.
func Base64Decoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	return io.ReadAll(decoder)
}

// Base 32 encoding.
func Base32Encoding(data []byte) []byte {
	var b32_buffer bytes.Buffer
	encoder := base32.NewEncoder(base32.StdEncoding, &b32_buffer)
	encoder.Write(data)
	encoder.Close()
	return b32_buffer.Bytes()
}

// Base 32 decoding.
func Base32Decoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := base32.NewDecoder(base32.StdEncoding, buf)
	return io.ReadAll(decoder)
}

// Hex encoding.
func Base16Encoding(data []byte) []byte {
	return []byte(hex.EncodeToString(data))
}

// Hex decoding.
func Base16Decoding(data []byte) ([]byte, error) {
	return hex.DecodeString(string(data))
}

// Gob encoding.
func GobEncoding(data []byte) []byte {
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(data)
	return buf.Bytes()
}

// Gob decoding.
func GobDecoding(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var decoded []byte
	err := decoder.Decode(&decoded)
	return decoded, err
}

func KeyEncoder(key *[32]byte) []byte {
	b_key := key[:]
	return GobEncoding(b_key)
}

func KeyDecoder(data []byte) (*[32]byte, error) {
	b_key, err := GobDecoding(data)
	if err != nil {
		return nil, err
	}
	var key [32]byte
	copy(key[:], b_key)
	return &key, nil
}

func Compress(data []byte) ([]byte, error) {
	// Create a new buffer
	buffer := new(bytes.Buffer)

	// Create a new gzip writer
	writer := gzip.NewWriter(buffer)

	// Write the data to the writer
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	// Create a new buffer
	buffer := bytes.NewBuffer(data)

	// Create a new gzip reader
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	// Read the data from the reader
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Close the reader
	err = reader.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
