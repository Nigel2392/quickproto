package quickproto

import (
	"bytes"
	"net"

	"github.com/Nigel2392/simplecrypto/aes"
)

func ReadConn(conn net.Conn, conf *Config, aes_key *[32]byte) (*Message, error) {
	msg := NewMessage(conf.Delimiter, conf.UseEncoding, conf.Enc_func, conf.Dec_func)
	buf := make([]byte, conf.BufSize)
	var data []byte
	// read until ending delimiter is found
	for !bytes.Contains(data, msg.EndingDelimiter()) {
		// read data from connection
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf[:n]...)
		// flush buffer
		if !bytes.Contains(data, msg.EndingDelimiter()) {
			buf = make([]byte, conf.BufSize)
		}
	}
	// decrypt data if needed
	if aes_key != nil {
		var err error
		data = bytes.TrimSuffix(data, msg.EndingDelimiter())
		data, err = aes.Decrypt(data, aes_key)
		if err != nil {
			return nil, err
		}
		data = append(data, msg.EndingDelimiter()...)
	}
	msg.Data = data
	return msg.Parse()
}

func WriteConn(conn net.Conn, msg *Message, aes_key *[32]byte) error {
	// Write data to connection
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
		send.Data = append(send.Data, msg.EndingDelimiter()...)
	}
	_, err = conn.Write(send.Data)
	if err != nil {
		return err
	}
	return nil
}
