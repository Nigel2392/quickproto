package quickproto

import (
	"bytes"
	"encoding/base64"
	"errors"
)

// A Message is a protocol message.
type Message struct {
	Data       []byte
	Delimiter  []byte
	Headers    map[string]string
	Body       []byte
	Use_Base64 bool
	// Parsed    bool
	// Generated bool
}

// NewMessage creates a new Message.
func NewMessage(delimiter []byte, use_b64 bool) *Message {
	return &Message{
		Data:       []byte{},
		Delimiter:  delimiter,
		Headers:    make(map[string]string),
		Body:       []byte{},
		Use_Base64: use_b64,
		// Parsed:    false,
		// Generated: false,
	}
}

// parses protocol messages.
// Header is a map of key/value pairs.
// Body is a base64 encoded byte slice.
func (m *Message) Parse() (*Message, error) {
	// Check if the message has already been parsed
	// if m.Parsed || m.Generated {
	// 	return m, errors.New("message has already been parsed")
	// }
	header_delimiter := append(m.Delimiter, m.Delimiter...)
	body_delimiter := append(header_delimiter, header_delimiter...)
	// Split data into headers and body
	data := bytes.SplitN(m.Data, body_delimiter, 2)
	if len(data) != 2 {
		return nil, errors.New("invalid message sent")
	}
	// Get headers from m.Data
	headers := bytes.Split(data[0], header_delimiter)
	for _, header := range headers {
		// Split header into key and value
		head := bytes.SplitN(header, m.Delimiter, 2)
		m.Headers[string(head[0])] = string(head[1])
	}
	// Get body from m.Data
	// Decode base64 encoded body
	var body []byte
	var err error
	if m.Use_Base64 {
		b64 := data[1]
		body, err = base64.StdEncoding.DecodeString(string(b64))
		if err != nil {
			return nil, errors.New("invalid base64 encoded body")
		}
	} else {
		body = data[1]
	}
	m.Body = body
	// m.Parsed = true
	return m, nil
}

// creates a protocol message.
// Header is a map of key/value pairs.
// Body is a base64 encoded byte slice.
func (m *Message) Generate() (*Message, error) {
	// Check if the message has already been generated
	// if m.Generated || m.Parsed {
	// 	return m, errors.New("message has already been generated")
	// }
	var buffer bytes.Buffer
	header_delimiter := append(m.Delimiter, m.Delimiter...)
	// Create headers
	for key, value := range m.Headers {
		// Create buffer for length of current header line
		headerline := make([]byte, len(key)+len(value)+len(m.Delimiter)+len(header_delimiter))
		// Copy the header into the headerline
		copy(headerline, []byte(key))
		copy(headerline[len(key):], m.Delimiter)
		copy(headerline[len(key)+len(m.Delimiter):], []byte(value))
		copy(headerline[len(key)+len(m.Delimiter)+len(value):], header_delimiter)
		// Append the headerline to the buffer
		buffer.Write(headerline)
	}
	// Create body
	buffer.Write(header_delimiter)
	if m.Use_Base64 {
		b64 := base64.StdEncoding.EncodeToString(m.Body)
		buffer.WriteString(b64)
	} else {
		buffer.Write(m.Body)
	}
	m.Data = buffer.Bytes()
	// m.Generated = true
	return m, nil
}

// Get content length of the message.
func (m *Message) ContentLength() int {
	return len(m.Data)
}
