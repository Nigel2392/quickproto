package quickproto

import (
	"bytes"
	"net"
)

func ReadConn(conn net.Conn, delimiter []byte, use_b64 bool, buf_size int) (*Message, error) {
	msg := NewMessage(delimiter, use_b64)
	buf := make([]byte, buf_size)
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
			buf = make([]byte, buf_size)
		}
	}
	msg.Data = data
	return msg.Parse()
}

func WriteConn(conn net.Conn, msg *Message) error {
	// Write data to connection
	send, err := msg.Generate()
	if err != nil {
		return err
	}
	_, err = conn.Write(send.Data)
	if err != nil {
		return err
	}
	return nil
}
