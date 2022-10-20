package quickproto

import (
	"bytes"
	"net"
)

func ReadConn(conn net.Conn, delimiter []byte, use_b64 bool, buf_size int) (*Message, error) {
	// Read data from connection, read until ending delimiter
	msg := NewMessage(delimiter, use_b64)
	buf := make([]byte, buf_size)
	var data []byte
	for !bytes.Contains(data, msg.EndingDelimiter()) {
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf[:n]...)
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
