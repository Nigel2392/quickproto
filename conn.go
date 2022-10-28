package quickproto

import (
	"bytes"
	"net"
	"strconv"

	"github.com/Nigel2392/simplecrypto/aes"
)

// Convenience function to craft an address from an IP and a port.
// Port could be string or int, IP must be string.
func CraftAddr(ip string, port any) string {
	switch port := port.(type) {
	case int:
		return ip + ":" + strconv.Itoa(port)
	case string:
		return ip + ":" + port
	default:
		panic("invalid port type provided")
	}
}

// ReadConn reads a message from a connection.
func ReadConn(conn net.Conn, conf *Config, aes_key *[32]byte, compress bool) (*Message, error) {
	msg := conf.NewMessage()
	buf := make([]byte, conf.BufSize)
	var data []byte
	// read until ending delimiter is found.
	for !bytes.Contains(data, msg.EndingDelimiter()) {
		// read data from connection.
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf[:n]...)
		// flush buffer.
	}
	// decrypt data if needed.
	if compress {
		data = bytes.TrimSuffix(data, msg.EndingDelimiter())
		var err error
		data, err = GZIPdecompress(data)
		if err != nil {
			return nil, err
		}
	}
	if aes_key != nil {
		var err error
		if !compress {
			data = bytes.TrimSuffix(data, msg.EndingDelimiter())
		}
		if data, err = aes.Decrypt(data, aes_key); err != nil {
			return nil, err
		}
		data = append(data, msg.EndingDelimiter()...)
	}
	msg.Data = data
	return msg.Parse()
}

// WriteConn writes a message to a connection and encrypts it if needed.
func WriteConn(conn net.Conn, msg *Message, aes_key *[32]byte, compress bool) error {
	// Write data to connection.
	send, err := msg.Generate()
	if err != nil {
		return err
	}
	if aes_key != nil {
		send.Data = bytes.TrimSuffix(send.Data, msg.EndingDelimiter())
		send.Data, err = aes.Encrypt(send.Data, aes_key)
		if err != nil {
			return err
		}
		if !compress {
			send.Data = append(send.Data, msg.EndingDelimiter()...)
		}
	}
	if compress {
		send.Data, err = GZIPcompress(send.Data)
		if err != nil {
			return err
		}
		send.Data = append(send.Data, msg.EndingDelimiter()...)
	}
	_, err = conn.Write(send.Data)
	if err != nil {
		return err
	}
	return nil
}
